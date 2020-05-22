// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package parser

//
//func TestCommandAUTHCONFIG(t *testing.T) {
//	tests := []commandTest{
//		{
//			name: "AUTHCONFIG/all params",
//			source: func() string {
//				return "AUTHCONFIG username:test-user private-key:/a/b/c"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAuthConfig]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAuthConfig)
//				}
//				authCmd, ok := cmds[0].(*AuthConfigCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if authCmd.GetUsername() != "test-user" {
//					return fmt.Errorf("Unexpected username %s", authCmd.GetUsername())
//				}
//				if authCmd.GetPrivateKey() != "/a/b/c" {
//					return fmt.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AUTHCONFIG - quoted params",
//			source: func() string {
//				return "AUTHCONFIG username:test-user private-key:'/a/b/c'"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAuthConfig]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAuthConfig)
//				}
//				authCmd, ok := cmds[0].(*AuthConfigCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if authCmd.GetUsername() != "test-user" {
//					return fmt.Errorf("Unexpected username %s", authCmd.GetUsername())
//				}
//				if authCmd.GetPrivateKey() != "/a/b/c" {
//					return fmt.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AUTHCONFIG with only private-key",
//			source: func() string {
//				return "AUTHCONFIG private-key:/a/b/c"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAuthConfig]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAuthConfig)
//				}
//				authCmd, ok := cmds[0].(*AuthConfigCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if authCmd.GetUsername() != "" {
//					return fmt.Errorf("Unexpected username %s", authCmd.GetUsername())
//				}
//				if authCmd.GetPrivateKey() != "/a/b/c" {
//					return fmt.Errorf("Unexpected privateKey %s", authCmd.GetPrivateKey())
//				}
//				return nil
//			},
//		},
//		{
//			name: "AUTHCONFIG - with var expansion",
//			source: func() string {
//				os.Setenv("fookey", "/a/b/c")
//				return "AUTHCONFIG username:${USER} private-key:$fookey"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAuthConfig]
//				authCmd := cmds[0].(*AuthConfigCommand)
//				if authCmd.GetUsername() != ExpandEnv("$USER") {
//					return fmt.Errorf("Unexpected username %s", authCmd.GetUsername())
//				}
//				if authCmd.GetPrivateKey() != "/a/b/c" {
//					return fmt.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
//				}
//				return nil
//			},
//		},
//		{
//			name: "Multiple AUTHCONFIG provided",
//			source: func() string {
//				return "AUTHCONFIG private-key:/foo/bar\nAUTHCONFIG username:test-user"
//			},
//			script: func(s *Script) error {
//				return nil
//			},
//			shouldFail: true,
//		},
//		{
//			name: "AUTHCONFIG with bad args",
//			source: func() string {
//				return "SSHCONFIG bar private-key:buzz"
//			},
//			shouldFail: true,
//		},
//
//		{
//			name: "AUTHCONFIG - with embedded colon",
//			source: func() string {
//				return "AUTHCONFIG username:test-user private-key:'/a/:b/c'"
//			},
//			script: func(s *Script) error {
//				cmds := s.Preambles[CmdAuthConfig]
//				if len(cmds) != 1 {
//					return fmt.Errorf("Script missing preamble %s", CmdAuthConfig)
//				}
//				authCmd, ok := cmds[0].(*AuthConfigCommand)
//				if !ok {
//					return fmt.Errorf("Unexpected type %T in script", cmds[0])
//				}
//				if authCmd.GetUsername() != "test-user" {
//					return fmt.Errorf("Unexpected username %s", authCmd.GetUsername())
//				}
//				if authCmd.GetPrivateKey() != "/a/:b/c" {
//					return fmt.Errorf("Unexpected private-key %s", authCmd.GetPrivateKey())
//				}
//				return nil
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			runCommandTest(t, test)
//		})
//	}
//}
