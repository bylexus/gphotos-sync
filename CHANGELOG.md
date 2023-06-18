# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Wishlist]

- credentials store encryption
- output folder configuration: folder level generation by template (e.g. "{year}/{month}")
- Filters: by Media Type
- Extract photo metadata from google api and embed them as EXIF data
- Extract GPS data: Unfortunately, it seems that those are NOT available through the Photos API. It _may_ be a possibility to
  use Google Drive: https://www.labnol.org/code/20059-image-exif-and-location
- Improve logging and error handling: Separate errors on stderr, use a logger
- Sync back: Upload changed / missing photos to Google Photos
- make "releases" infrastructure to publish pre-compiled binaries

## [Unreleased]

- bug: Fixing Date bug: 1st and last of day/month caused a wrong format error
- feature: by date range: It is now possible to filter by date range: `--date-range=2023-04-01:2023-05-15`

## [0.1.0] - 2023-06-17

This is the first release ever - the [Added] section below just lists the actual features.

### Added

- downloads Google Photo media items to a local folder
- creates a folder for each year
- supported files: Photos and Videos, as supported by Google Photos
- support for multiple download threads
- support for skip / override existing / newer files
- stores the google api tokens in a local credentials store
