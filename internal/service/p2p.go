package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/p2p"
	peer "github.com/libp2p/go-libp2p/core/peer"
	pstore "github.com/libp2p/go-libp2p/core/peerstore"
	protocol "github.com/libp2p/go-libp2p/core/protocol"
	pb "github.com/mohaijiang/computeshare-client/api/network/v1"
	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
	"strconv"
	"strings"
	"time"
)

const P2PProtoPrefix = "/x/"

var resolveTimeout = 10 * time.Second

type P2pService struct {
	pb.UnimplementedP2PServer
	n *core.IpfsNode
}

func NewP2pService(ipfsNode *core.IpfsNode) *P2pService {
	return &P2pService{
		n: ipfsNode,
	}
}

func (s *P2pService) CreateListen(ctx context.Context, req *pb.CreateListenRequest) (*pb.CreateListenReply, error) {
	protoOpt := req.GetProtocol()
	targetOpt := req.GetTargetAddress()

	proto := protocol.ID(protoOpt)

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		return nil, err
	}

	// port can't be 0
	if err := checkPort(target); err != nil {
		return nil, err
	}

	allowCustom := false
	reportPeerID := false

	if !allowCustom && !strings.HasPrefix(string(proto), P2PProtoPrefix) {
		return nil, errors.New("protocol name must be within '" + P2PProtoPrefix + "' namespace")
	}

	_, err = s.n.P2P.ForwardRemote(s.n.Context(), proto, target, reportPeerID)
	return &pb.CreateListenReply{}, err
}
func (s *P2pService) CreateForward(ctx context.Context, req *pb.CreateForwardRequest) (*pb.CreateForwardReply, error) {
	protoOpt := req.GetPortal()
	listenOpt := req.GetListenAddress()
	targetOpt := req.GetTargetAddress()

	proto := protocol.ID(protoOpt)

	listen, err := ma.NewMultiaddr(listenOpt)
	if err != nil {
		return nil, err
	}

	targets, err := parseIpfsAddr(targetOpt)
	if err != nil {
		return nil, err
	}

	allowCustom := false

	if !allowCustom && !strings.HasPrefix(string(proto), P2PProtoPrefix) {
		return nil, errors.New("protocol name must be within '" + P2PProtoPrefix + "' namespace")
	}

	err = forwardLocal(s.n.Context(), s.n.P2P, s.n.Peerstore, proto, listen, targets)
	return &pb.CreateForwardReply{}, err
}
func (s *P2pService) CloseListen(ctx context.Context, req *pb.CloseListenRequest) (*pb.CloseListenReply, error) {

	var proto protocol.ID
	if req.Protocol != nil {
		proto = protocol.ID(req.GetProtocol())
	}

	var target, listen ma.Multiaddr
	var err error

	if req.ListenAddress != nil {
		listen, err = ma.NewMultiaddr(req.GetListenAddress())
		if err != nil {
			return nil, err
		}
	}

	if req.TargetAddress != nil {
		target, err = ma.NewMultiaddr(req.GetTargetAddress())
		if err != nil {
			return nil, err
		}
	}

	match := func(listener p2p.Listener) bool {
		if proto != "" && proto != listener.Protocol() {
			return false
		}
		if listen != nil && !listen.Equal(listener.ListenAddress()) {
			return false
		}
		if target != nil && !target.Equal(listener.TargetAddress()) {
			return false
		}
		return true
	}

	done := s.n.P2P.ListenersLocal.Close(match)
	done += s.n.P2P.ListenersP2P.Close(match)
	return &pb.CloseListenReply{}, nil
}
func (s *P2pService) ListListen(ctx context.Context, req *pb.ListListenRequest) (*pb.ListListenReply, error) {
	output := &pb.ListListenReply{}

	s.n.P2P.ListenersLocal.Lock()
	for _, listener := range s.n.P2P.ListenersLocal.Listeners {
		output.Result = append(output.Result, &pb.ListenReply{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	s.n.P2P.ListenersLocal.Unlock()

	s.n.P2P.ListenersP2P.Lock()
	for _, listener := range s.n.P2P.ListenersP2P.Listeners {
		output.Result = append(output.Result, &pb.ListenReply{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	s.n.P2P.ListenersP2P.Unlock()
	return output, nil
}

// checkPort checks whether target multiaddr contains tcp or udp protocol
// and whether the port is equal to 0
func checkPort(target ma.Multiaddr) error {
	// get tcp or udp port from multiaddr
	getPort := func() (string, error) {
		sport, _ := target.ValueForProtocol(ma.P_TCP)
		if sport != "" {
			return sport, nil
		}

		sport, _ = target.ValueForProtocol(ma.P_UDP)
		if sport != "" {
			return sport, nil
		}
		return "", fmt.Errorf("address does not contain tcp or udp protocol")
	}

	sport, err := getPort()
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(sport)
	if err != nil {
		return err
	}

	if port == 0 {
		return fmt.Errorf("port can not be 0")
	}

	return nil
}

// parseIpfsAddr is a function that takes in addr string and return ipfsAddrs
func parseIpfsAddr(addr string) (*peer.AddrInfo, error) {
	multiaddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	pi, err := peer.AddrInfoFromP2pAddr(multiaddr)
	if err == nil {
		return pi, nil
	}

	// resolve multiaddr whose protocol is not ma.P_IPFS
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	defer cancel()
	addrs, err := madns.Resolve(ctx, multiaddr)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.New("fail to resolve the multiaddr:" + multiaddr.String())
	}
	var info peer.AddrInfo
	for _, addr := range addrs {
		taddr, id := peer.SplitAddr(addr)
		if id == "" {
			// not an ipfs addr, skipping.
			continue
		}
		switch info.ID {
		case "":
			info.ID = id
		case id:
		default:
			return nil, fmt.Errorf(
				"ambiguous multiaddr %s could refer to %s or %s",
				multiaddr,
				info.ID,
				id,
			)
		}
		info.Addrs = append(info.Addrs, taddr)
	}
	return &info, nil
}

// forwardLocal forwards local connections to a libp2p service
func forwardLocal(ctx context.Context, p *p2p.P2P, ps pstore.Peerstore, proto protocol.ID, bindAddr ma.Multiaddr, addr *peer.AddrInfo) error {
	ps.AddAddrs(addr.ID, addr.Addrs, pstore.TempAddrTTL)
	// TODO: return some info
	_, err := p.ForwardLocal(ctx, addr.ID, proto, bindAddr)
	return err
}
