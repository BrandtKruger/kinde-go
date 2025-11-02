#!/bin/bash
# Post-generation fix for OptString.Decode() null handling
# This script patches the generated code until ogen fixes the issue upstream

set -e

FILE="oas_json_gen.go"

echo "Applying OptString.Decode() null handling fix to $FILE..."

# Create a backup
cp "$FILE" "$FILE.backup"

# Use sed to replace the OptString.Decode method
# This is fragile but works until ogen is fixed
perl -i -0pe 's/func \(o \*OptString\) Decode\(d \*jx\.Decoder\) error \{\n\tif o == nil \{\n\t\treturn errors\.New\("invalid: unable to decode OptString to nil"\)\n\t\}\n\to\.Set = true\n\tv, err := d\.Str\(\)/func (o *OptString) Decode(d *jx.Decoder) error {\n\tif o == nil {\n\t\treturn errors.New("invalid: unable to decode OptString to nil")\n\t}\n\t\/\/ Check if the value is null and handle it gracefully\n\tif d.Next() == jx.Null {\n\t\tif err := d.Null(); err != nil {\n\t\t\treturn err\n\t\t}\n\t\t\/\/ For null values, treat as unset (not present)\n\t\to.Set = false\n\t\to.Value = ""\n\t\treturn nil\n\t}\n\t\/\/ Value is present and not null\n\to.Set = true\n\tv, err := d.Str()/g' "$FILE"

echo "Fix applied successfully!"
echo "Backup saved as $FILE.backup"

