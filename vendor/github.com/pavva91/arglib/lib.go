/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package arglib

import (
	"errors"
	"fmt"
	"strconv"
)

// var complexInteractionsLog = log.Logger{}("complexInteractions")
func ArgumentSizeVerification(args []string, numberOfArguments int) error {
	if len(args) != numberOfArguments {
		return errors.New("Incorrect number of arguments. Expecting " + string(numberOfArguments))
	}
	return nil
}
func ArgumentSizeLimitVerification(args []string, numberOfArguments int) error {
	if len(args) > numberOfArguments {
		return errors.New("Incorrect number of arguments. Expecting " + string(numberOfArguments))
	}
	return nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func SanitizeArguments(args []string) error {
	for i, arg := range args {
		fmt.Print(i)
		fmt.Println(arg)

		if len(arg) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		// TODO: Fix, on serviceEvaluation the id now is too long, after ID definitions
		// if len(arg) > 32 {
		// 	return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		// }
	}
	return nil
}
