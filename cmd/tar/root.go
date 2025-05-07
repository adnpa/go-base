package tar

import (
	"github.com/spf13/cobra"
)

var TarCmd = &cobra.Command{
	Use:   "tar",
	Short: "压缩当前目录所有文件",
	Long:  `压缩当前目录所有文件`,
}

func init() {
	TarCmd.AddCommand(allCmd)
}
