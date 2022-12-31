Feature: Testsuite test
  Quentin check whether their testsuite works properly.

  @linux @darwin @windows
  Scenario: Contains
    When executing "go help" succeeds
    Then stdout should contain
      """
      Go is a tool for managing Go source code.
      """

  @linux @darwin
  Scenario: Command which is not present
    When executing "foobar" fails
    Then exitcode should not equal "0"
    Then stderr should contain
      """
      command not found
      """

  @windows
  Scenario: Command which is not present
    When executing "foobar" succeed
    Then stderr should contain
      """
      'foobar' is not recognized
      """

  @linux @darwin @windows
  Scenario: Not Contains
    When executing "go help" succeeds
    Then stdout should not contain "Error"

  @linux @darwin @windows
  Scenario: Equals
    When executing "go help" succeeds
    Then exitcode should equal "0"

  @linux @darwin @windows
  Scenario: Not Equals
    When executing "go notexist" fails
    Then exitcode should not equal "0"

  @linux @darwin @windows
  Scenario: Matches
    When executing "go version" succeeds
    Then stdout should match "go1\.\d+\.\d+"

  @linux @darwin @windows
  Scenario: Not Matches
    When executing "go version" succeeds
    Then stdout should not match "Local \d\.\d OpenShift clusters"

  @linux @darwin @windows
  Scenario: Is Empty
    When executing "go version" succeeds
    Then stderr should be empty

  @linux @darwin @windows
  Scenario: Is Not Empty
    When executing "go version" succeeds
    Then stdout should not be empty

  @linux @darwin
  Scenario: Able to use defined variable
    When executing "VAR=$(echo 'hello')" succeeds
    And executing "echo $VAR" succeeds
    Then stdout should contain
      """
      hello
      """

  @windows
  Scenario: Able to use defined variable
    When executing "$Env:POD = $(echo 'hello')" succeeds
    And executing "echo $Env:POD" succeeds
    Then stdout should contain
      """
      hello
      """

  @linux @darwin
  Scenario: Scenario Variables
    When setting scenario variable "VAR" to the stdout from executing "go version"
    And executing "echo $(VAR)" succeeds
    Then stdout should contain "go version"

  @windows
  Scenario: Scenario Variables
    When setting scenario variable "VAR" to the stdout from executing "go version"
    And executing "echo $(VAR)" succeeds
    Then stdout should contain "version"

  @linux @darwin @windows
  Scenario: Create Directory and Files
    When creating directory "newdir" succeeds
    And creating file "newdir/newfile" succeeds
    And file from "https://google.com" is downloaded into location "newdir"

  @linux @darwin @windows
  Scenario: File Content Checks
    Given creating directory "newdir" succeeds
    And creating file "newdir/newfile" succeeds
    When writing text "192.168.15.17" to file "newdir/newfile" succeeds
    Then content of file "newdir/newfile" should contain "168"
    And content of file "newdir/newfile" should not contain "512"
    And content of file "newdir/newfile" should equal "192.168.15.17"
    And content of file "newdir/newfile" should not equal "192.168.15.512"
    And content of file "newdir/newfile" should match "192\.168\.\d+\.\d+"
    And content of file "newdir/newfile" should not match "192\.168\.\s+\.\d+"
    And content of file "newdir/newfile" is valid "IP"

  @linux @darwin @windows
  Scenario: Delete file
    When deleting file "newdir/newfile" succeeds
    Then file "newdir/newfile" should not exist

  @linux @darwin @windows
  Scenario: Delete directory
    When deleting directory "newdir" succeeds
    Then directory "newdir" should not exist

   # Config

  @linux @darwin @windows
  Scenario Outline: Verify key exists in JSON/YAML config file
    Given file "<filename>" exists
    When  "<format>" config file "<filename>" contains key "<property>"
    And   "<format>" config file "<filename>" does not contain key "<nonproperty>"
    And   "<format>" config file "<filename>" does not contain key "whocareswhatkey"
    Then  "<format>" config file "<filename>" contains key "<property>" with value matching "<goodvalue>"
    And   "<format>" config file "<filename>" does not contain key "<property>" with value matching "<badvalue>"
    And   "<format>" config file "<filename>" does not contain key "<property>" with value matching "whocareswhatvalue"

    Examples: Config files contain keys
      | format | filename                       | property          | nonproperty | goodvalue  | badvalue |
      | JSON   | ../../testdata/testconfig.json | author.login      | nonauthor   | alice      | bob      |
      | JSON   | ../../testdata/testconfig.json | tag_name          | nontag      | v1.3.1     | v2.2.2   |
      | JSON   | ../../testdata/testconfig.json | tag_name          | nontag      | v1\.\d+\.1 | v2.2.2   |
      | YAML   | ../../testdata/testconfig.yml  | workflows.version | nonjobs     | zero       | one      |
      | YAML   | ../../testdata/testconfig.yml  | version           | nonversion  | 2          | two      |
      | YAML   | ../../testdata/testconfig.yml  | version           | nonversion  | \d+        | two      |

  @linux @darwin @windows
  Scenario: Verify that key matching value exists in JSON config file
    When file "../../testdata/testconfig.json" exists
      # JSON package uses map[string]interface{} and []interface{} values to store various JSON objects
      # and arrays. It unmarshalls into 4 categories: bool, float64 (all numbers), string, nil.
      # Hence, id=4412902 needs to be given in float64 format to get a match:
    Then "JSON" config file "../../testdata/testconfig.json" contains key "id" with value matching "4.412902e\+06"
