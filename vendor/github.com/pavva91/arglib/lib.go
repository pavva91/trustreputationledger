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
	"strconv"
)

func ArgumentSizeVerification(args []string, numberOfArguments int) error {
	if len(args) != numberOfArguments {
		return errors.New("Incorrect number of arguments. Expecting " + string(numberOfArguments))
	}
	return nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func SanitizeArguments(args []string) error {
	for i, arg := range args {
		if len(arg) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(arg) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}
