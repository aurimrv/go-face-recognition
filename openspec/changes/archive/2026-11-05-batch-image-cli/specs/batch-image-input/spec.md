## ADDED Requirements

### Requirement: Default single-image analysis
When invoked with no positional path arguments, the system SHALL analyze `images/unknown.jpg` and write the annotated image to `images/result.jpg`.

#### Scenario: No arguments uses historic default paths
- **WHEN** the program is started with no path arguments
- **THEN** the system analyzes `images/unknown.jpg` and saves the result as `images/result.jpg`

#### Scenario: Default mode failure is fatal
- **WHEN** the program is started with no path arguments and classification or drawing fails for the default image
- **THEN** the program exits with a non-zero status and does not continue

### Requirement: Explicit file arguments
The system SHALL treat each positional argument that names an existing regular file as a single analysis job and write the annotated image beside the source as `result_<stem>.jpg`, where `<stem>` is the source filename without its extension.

#### Scenario: Multiple file arguments
- **WHEN** the program is started with arguments `images/a.jpg` and `images/b.jpg`
- **THEN** the system analyzes both files and writes `images/result_a.jpg` and `images/result_b.jpg`

#### Scenario: Nested file path keeps output beside source
- **WHEN** the program is started with argument `images/Kate-Winslet/unknown-01.jpg`
- **THEN** the system writes `images/Kate-Winslet/result_unknown-01.jpg`

### Requirement: Directory arguments expand recursively
The system SHALL treat each positional argument that names a directory as a recursive walk of that directory and its subdirectories, producing one analysis job per eligible image file.

#### Scenario: Album directory produces colocated results
- **WHEN** the program is started with a directory argument containing JPEG images in nested subfolders
- **THEN** each eligible image is analyzed and a `result_<stem>.jpg` is written in the same directory as that source image

#### Scenario: Mixed files and directories
- **WHEN** the program is started with both file and directory path arguments
- **THEN** the system flattens all resolved jobs into one list and processes each job

### Requirement: Directory walk eligibility filters
When expanding a directory, the system SHALL include only files whose extension is `.jpg` or `.jpeg` (case-insensitive), and SHALL exclude files whose basename is `result.jpg` or matches the pattern `result_*.jpg`.

#### Scenario: Non-JPEG files are skipped during walk
- **WHEN** a walked directory contains `.avif` or other non-JPEG files alongside `.jpg` files
- **THEN** only the JPEG files become analysis jobs

#### Scenario: Prior result images are skipped during walk
- **WHEN** a walked directory contains `photo.jpg` and `result_photo.jpg`
- **THEN** only `photo.jpg` becomes an analysis job

### Requirement: Invalid or empty job resolution
The system SHALL fail fatally if a positional path argument does not exist or is neither a regular file nor a directory, and SHALL fail fatally if path resolution yields zero analysis jobs.

#### Scenario: Missing path argument
- **WHEN** the program is started with a path that does not exist
- **THEN** the program exits with a non-zero status before analyzing any images

#### Scenario: Directory with no eligible images
- **WHEN** the program is started with a directory that contains no eligible JPEG inputs after filters
- **THEN** the program exits with a non-zero status

### Requirement: Shared setup and sequential processing
The system SHALL load persons and initialize the face recognizer once per process, then process resolved analysis jobs sequentially (classify, annotate, save) in the resolved order.

#### Scenario: One-time recognizer setup for batch
- **WHEN** multiple analysis jobs are resolved from the arguments
- **THEN** person loading and recognizer initialization run once before the first job

### Requirement: Batch error handling
When more than one analysis job is resolved, or when any explicit path argument was provided, the system SHALL log a per-job error and continue with remaining jobs if classification, drawing, or saving fails for one image, and SHALL exit with a non-zero status if any job failed.

#### Scenario: One failure does not abort remaining jobs
- **WHEN** two or more jobs are queued and the first job fails during classification
- **THEN** the system reports the error for that job and continues processing the remaining jobs

#### Scenario: Any job failure yields non-zero exit
- **WHEN** at least one job in a multi-job run fails
- **THEN** after all jobs are attempted the program exits with a non-zero status
