// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package script

import (
	"fmt"
	"os"
	"testing"
)

func TestCommandAS(t *testing.T) {
	tests := []commandTest{
		{
			name: "AS/unquoted",
			command: func() Command {
				cmd, err := NewAsCommand(0, "userid:foo groupid:bar")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},

			test: func(t *testing.T, cmd Command) {
				asCmd, ok := cmd.(*AsCommand)
				if !ok {
					t.Fatalf("Unexpected type %T in script", cmd)
				}
				if asCmd.GetUserId() != "foo" {
					t.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
				}
				if asCmd.GetGroupId() != "bar" {
					t.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
				}
			},
		},
		{
			name: "AS/quoted",
			command: func() Command {
				cmd, err :=  NewAsCommand(0,`userid:"foo" groupid:bar`)
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, cmd Command) {
				asCmd, ok := cmd.(*AsCommand)
				if !ok {
					t.Fatalf("Unexpected type %T in script", cmd)
				}
				if asCmd.GetUserId() != "foo" {
					t.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
				}
				if asCmd.GetGroupId() != "bar" {
					t.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
				}
			},
		},
		{
			name: "AS/userid only",
			command: func() Command{
				cmd, err := NewAsCommand(0, "userid:foo")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, cmd Command) {
				asCmd, ok := cmd.(*AsCommand)
				if !ok {
					t.Fatalf("Unexpected type %T in script", cmd)
				}
				if asCmd.GetUserId() != "foo" {
					t.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
				}
				if asCmd.GetGroupId() != fmt.Sprintf("%d", os.Getgid()) {
					t.Errorf("Unexpected AS groupid %s", asCmd.GetGroupId())
				}
			},
		},

		{
			name: "AS/var expansion",
			command: func() Command {
				cmd, err := NewAsCommand(0, "userid:$USER groupid:$foogid")
				if err != nil {
					t.Fatal(err)
				}
				return cmd
			},
			test: func(t *testing.T, cmd Command) {
				os.Setenv("foogid", "barid")
				asCmd := cmd.(*AsCommand)
				if asCmd.GetUserId() != ExpandEnv("$USER") {
					t.Errorf("Unexpected AS userid %s", asCmd.GetUserId())
				}
				if asCmd.GetGroupId() != "barid" {
					t.Errorf("Unexpected AS groupid %s", asCmd.GetUserId())
				}
			},
		},
		{
			name: "AS/unsupported args",
			command: func() Command {
				cmd, err := NewAsCommand(0, "foo:bar fuzz:buzz")
				if err == nil {
					t.Fatal("Expected failure, but err is nil")
				}
				return cmd
			},
			test:func(t *testing.T, cmd Command){
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runCommandTest(t, test)
		})
	}
}
