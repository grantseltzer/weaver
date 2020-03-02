To generate new smoke test test_data:

- As functions are added to cmd/tester/main.go:
    - make sure they're called in main
    - add them to tests/test_data/functions_file.txt

- Run weaver, and then the tester. 
- Pipe output to a file. jq -s -c sorted by function name, saved to tester_output_sorted_slurped.json