package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

func main() {
    root := &cobra.Command{Use: "dockyard"}
    root.AddCommand(cmdProjects(), cmdDeploy())
    if err := root.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func cmdProjects() *cobra.Command {
    var name string
    cmd := &cobra.Command{Use: "projects", Short: "Manage projects"}
    create := &cobra.Command{Use: "create", RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("created project: %s\n", name) // { SPECULATION }
        return nil
    }}
    create.Flags().StringVar(&name, "name", "", "project name")
    _ = create.MarkFlagRequired("name")
    cmd.AddCommand(create)
    return cmd
}

func cmdDeploy() *cobra.Command {
    var id string
    cmd := &cobra.Command{Use: "deploy", RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("deploy triggered for %s\n", id)
        return nil
    }}
    cmd.Flags().StringVar(&id, "id", "", "project id")
    _ = cmd.MarkFlagRequired("id")
    return cmd
}
