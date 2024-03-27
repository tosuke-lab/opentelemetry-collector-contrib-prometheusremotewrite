#!/bin/bash

# Copyright The OpenTelemetry Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

GOJUNITREPORTCMD=${GOJUNIT:-go-junit-report}

TESTS_DIR=${TESTS_DIR:-tests}

cd ${TESTS_DIR}

SED="sed"

PASS_COLOR=$(printf "\033[32mPASS\033[0m")
FAIL_COLOR=$(printf "\033[31mFAIL\033[0m")
TEST_COLORIZE="${SED} 's/PASS/${PASS_COLOR}/' | ${SED} 's/FAIL/${FAIL_COLOR}/'"

mkdir -p results/junit

RUN_TESTBED=1 go test -v ${TEST_ARGS} 2>&1 | tee results/testoutput.log | bash -c "${TEST_COLORIZE}"

testStatus=${PIPESTATUS[0]}

${GOJUNITREPORTCMD} < results/testoutput.log > results/junit/results.xml

bash -c "cat results/TESTRESULTS.md | ${TEST_COLORIZE}"

exit ${testStatus}
