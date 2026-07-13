// Command torchbearing 运行一个完全确定性的 MS1 成功场景演示。
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spojchil/torchbearing/internal/bootstrap"
	"github.com/spojchil/torchbearing/internal/contracts"
	"github.com/spojchil/torchbearing/mocks/deterministic"
)

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(output io.Writer) error {
	gateway, err := bootstrap.NewMS1Gateway(deterministic.ScenarioSuccess)
	if err != nil {
		return err
	}

	response, failure := gateway.Analyze(context.Background(), contracts.AnalysisRequest{
		Text: "查看 checkout 服务过去 30 分钟的请求速率",
		Scope: contracts.AnalysisScope{
			DatasourceUID: "prometheus-mock",
			TimeRange: contracts.TimeRange{
				From: "now-30m",
				To:   "now",
			},
		},
	})
	if failure != nil {
		return fmt.Errorf("%s: %s", failure.Code, failure.Message)
	}
	return json.NewEncoder(output).Encode(response)
}
