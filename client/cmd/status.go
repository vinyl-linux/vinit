/*
Copyright © 2022 James Condron <james@zero-internet.org.uk>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	vinit "github.com/vinyl-linux/vinit/dispatcher"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status a service",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		c, err := newClient(socketAddr)
		if err != nil {
			return
		}

		var status *vinit.ServiceStatus

		if len(args) == 0 {
			// systemstatus
			var ss []*vinit.ServiceStatus

			ss, err = c.systemStatus()
			if err != nil {
				return
			}

			for _, status = range ss {
				fmt.Println(fmtStatus(status))
			}

			return
		}

		// specific service
		status, err = c.status(args[0])
		if err != nil {
			return
		}

		fmt.Println(fmtStatus(status))

		return
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func fmtStatus(s *vinit.ServiceStatus) string {
	return fmt.Sprintf("%s: %s\n%s %s\n%s",
		s.Svc.Name, runningStr(s.Running, s.Pid),
		startStr(s.StartTime.AsTime()), endStr(s.EndTime.AsTime()),
		completionDetails(s),
	)
}

func runningStr(b bool, pid uint32) string {
	if b {
		return color.HiGreenString("running") + fmt.Sprintf(" (pid: %d)", int(pid))
	}

	return color.HiBlackString("not running")
}

func startStr(t time.Time) string {
	return fmt.Sprintf("started at %s", t)
}

func endStr(t time.Time) string {
	if t.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return ""
	}

	return fmt.Sprintf("completed at %s", t)
}

func completionDetails(s *vinit.ServiceStatus) string {
	sb := new(strings.Builder)

	if s.Success || s.ExitStatus != 0 {
		sb.WriteString("last exit status " + fmt.Sprint(s.ExitStatus) + "\n")
	}

	if s.Error != "" {
		sb.WriteString(s.Error + "\n")
	}

	return sb.String()
}
