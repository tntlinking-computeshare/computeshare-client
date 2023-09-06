package service

import (
	pb "computeshare-client/api/compute/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/kubo/core"
	"github.com/ipfs/kubo/core/coreapi"
	"io"
	"os"
	"os/exec"
	path2 "path"
)

type ComputepowerService struct {
	pb.UnimplementedComputepowerServer
	n   *core.IpfsNode
	log *log.Helper
}

func NewComputepowerService(ipfsNode *core.IpfsNode, logger log.Logger) *ComputepowerService {
	return &ComputepowerService{
		n:   ipfsNode,
		log: log.NewHelper(logger),
	}
}

func (s *ComputepowerService) RunPythonPackage(ctx context.Context, req *pb.RunPythonPackageRequest) (*pb.RunPythonPackageReply, error) {
	api, err := coreapi.NewCoreAPI(s.n, options.Api.FetchBlocks(true))
	if err != nil {
		return nil, err
	}
	file, err := api.Unixfs().Get(context.Background(), path.New(req.GetCid()))
	if err != nil {
		return nil, err
	}

	size, err := file.Size()
	if err != nil {
		return nil, err
	}

	s.log.Info("file size: ", size)

	reader := files.ToFile(file)
	bytes, err := io.ReadAll(reader)
	err = os.WriteFile(path2.Join("/tmp", req.GetCid()), bytes, 0644)

	if err != nil {
		return nil, err
	}

	s.log.Info(string(bytes))

	// 执行python 命令

	return &pb.RunPythonPackageReply{}, nil
}
func (s *ComputepowerService) RunBenchmarks(ctx context.Context, req *pb.RunBenchmarksRequest) (*pb.RunBenchmarksReply, error) {
	cmd := exec.Command("sysbench", "cpu", "--cpu-max-prime=20000", "--threads=4", "run")
	output, err := cmd.CombinedOutput()
	return &pb.RunBenchmarksReply{
		Output: string(output),
	}, err
}
