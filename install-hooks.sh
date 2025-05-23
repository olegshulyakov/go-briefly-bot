#!/bin/sh

# Define the source and destination directories
SOURCE_DIR="hooks"
DEST_DIR=".git/hooks"

# Check if the source directory exists
if [ ! -d "$SOURCE_DIR" ]; then
  echo "Source directory '$SOURCE_DIR' does not exist."
  exit 1
fi

# Iterate over each file in the source directory
for FILE in "$SOURCE_DIR"/*; do
  # Get the base name of the file
  BASENAME=$(basename "$FILE")

  # Copy the file to the destination directory
  cp "$FILE" "$DEST_DIR/$BASENAME"

  # Make the copied file executable
  chmod +x "$DEST_DIR/$BASENAME"

  echo "Installed hook: $BASENAME"
done

echo "All hooks installed successfully."