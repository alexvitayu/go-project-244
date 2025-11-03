package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "gendiff",
		Usage:                  "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage: `output formats: 
				-f=stylish
				-f=plain
				-f=json`},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var path1 string
			var path2 string
			if cmd.Args().Len() > 0 {
				path1 = cmd.Args().Get(0)
				path2 = cmd.Args().Get(1)
			}
			format := cmd.String("format")
			result, err := code.GenDiff(path1, path2, format)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(result)
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
