# Synopsis  
This branch is to enable PostgreSQL as a database for PhotoPrism.    
The SQL for SQLite and MariaDB is not compatible with PostgreSQL (or each other) sometimes.  
This leads to failures.  

# Testing Status
The following shows the tests that were failing as at 2025-02-26.  As they are fixed, the status of the test will be updated.  
```
FAIL | AddPhotosToAlbum (0.01s)
FAIL |   AddPhotosToAlbum/AddMultiplePhotos (0.00s)
FAIL |   AddPhotosToAlbum/AddSinglePhoto (0.00s)
FAIL |   AddPhotosToAlbum/AddPhotoFromReview (0.00s)
FAIL | RemovePhotosFromAlbum (0.01s)
FAIL | BatchPhotosPrivate (0.01s)
FAIL |   BatchPhotosPrivate/Success (0.01s)
FAIL | BatchPhotosApprove (0.01s)
FAIL |   BatchPhotosApprove/Success (0.01s)
FAIL | GetFace (0.00s)
FAIL |   GetFace/Success (0.00s)
FAIL |   GetFace/Lowercase (0.00s)
FAIL | UpdateFace (0.00s)
FAIL |   UpdateFace/Success (0.00s)
FAIL | GetMomentsTime (0.00s)
FAIL |   GetMomentsTime/get_moments_time (0.00s)
FAIL | PhotoPrimary (0.00s)
FAIL |   PhotoPrimary/Success (0.00s)
FAIL | Zip (0.00s)
FAIL |   Zip/Download (0.00s)
FAIL | 	github.com/photoprism/photoprism/internal/api	56.180s
FAIL | ClientsListCommand (0.04s)
FAIL |   ClientsListCommand/Monitoring (0.01s)
FAIL |   ClientsListCommand/CSV (0.01s)
FAIL | UsersCommand (11.87s)
FAIL |   UsersCommand/AddModifyAndRemoveJohn (11.87s)
FAIL | 	github.com/photoprism/photoprism/internal/commands	147.531s
     | time="2025-02-26T04:50:00Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/entity	604.136s
FAIL | DialectMysql (0.07s)
FAIL |   DialectMysql/ValidMigration (0.04s)
FAIL |   DialectMysql/InvalidDataUpgrade (0.03s)
FAIL | 	github.com/photoprism/photoprism/internal/entity/dbtest	205.053s
FAIL | DialectMysql (0.05s)
FAIL | 	github.com/photoprism/photoprism/internal/entity/migrate	0.171s
PASS | UpdateAlbumDefaultCovers (0.00s)
PASS | UpdateAlbumFolderCovers (0.00s)
PASS | UpdateAlbumMonthCovers (0.00s)
PASS | UpdateAlbumCovers (0.00s)
PASS | UpdateLabelCovers (0.00s)
PASS | UpdateSubjectCovers (0.00s)
PASS | UpdateCovers (0.00s)
PASS | FileSelection (0.01s)
PASS |   FileSelection/DownloadSelectionRawSidecarPrivate (0.00s)
PASS |   FileSelection/DownloadSelectionRawOriginals (0.00s)
PASS |   FileSelection/ShareSelectionOriginals (0.00s)
PASS |   FileSelection/ShareSelectionPrimary (0.00s)
PASS |   FileSelection/ShareAlbums (0.00s)
PASS |   FileSelection/ShareMonths (0.00s)
PASS |   FileSelection/ShareFoldersOriginals (0.00s)
PASS |   FileSelection/ShareFolders (0.00s)
PASS |   FileSelection/ShareStatesOriginals (0.00s)
PASS |   FileSelection/ShareStates (0.00s)
PASS | SetDownloadFileID (0.00s)
PASS |   SetDownloadFileID/Success (0.00s)
PASS | FilesByUID (0.00s)
PASS |   FilesByUID/Negative_limit_with_offset (0.00s)
PASS | SetPhotoPrimary (0.00s)
PASS |   SetPhotoPrimary/Success (0.00s)
PASS |   SetPhotoPrimary/no_file_uid (0.00s)
PASS | AlbumFolders (0.00s)
PASS |   AlbumFolders/root (0.00s)
PASS | MomentsTime (0.00s)
PASS |   MomentsTime/PublicOnly (0.00s)
PASS |   MomentsTime/IncludePrivate (0.00s)
PASS | MomentsCountries (0.00s)
PASS |   MomentsCountries/PublicOnly (0.00s)
PASS |   MomentsCountries/IncludePrivate (0.00s)
PASS | MomentsStates (0.00s)
PASS |   MomentsStates/PublicOnly (0.00s)
PASS |   MomentsStates/IncludePrivate (0.00s)
PASS | MomentsCategories (0.00s)
PASS |   MomentsCategories/PublicOnly (0.00s)
PASS |   MomentsCategories/IncludePrivate (0.00s)
PASS | PhotoSelection (0.00s)
PASS |   PhotoSelection/photos_selected (0.00s)
PASS |   PhotoSelection/FindAlbums (0.00s)
PASS |   PhotoSelection/FindMonths (0.00s)
PASS |   PhotoSelection/FindFolders (0.00s)
PASS |   PhotoSelection/FindStates (0.00s)
PASS | FixPrimaries (0.00s)
PASS |   FixPrimaries/Success (0.00s)
FAIL | 	github.com/photoprism/photoprism/internal/entity/query	100.671s
FAIL | Albums (0.00s)
FAIL |   Albums/search_with_string (0.00s)
************************************************************
Panic error caused all subsequent tests to not be run.
Panic was caused because the test above returned no results.
Query expects the DBMS to be be case insensitive.
PostgreSQL is not by default.
************************************************************
FAIL | 	github.com/photoprism/photoprism/internal/entity/search	88.476s
     | time="2025-02-26T04:52:41Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/photoprism	604.100s
     | time="2025-02-26T04:53:44Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/photoprism/backup	604.186s
     | time="2025-02-26T04:53:50Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/photoprism/get	604.205s
     | time="2025-02-26T05:00:06Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/server	604.195s
     | time="2025-02-26T05:02:47Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/server/wellknown	604.192s
     | time="2025-02-26T05:03:58Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/thumb/avatar	604.288s
     | time="2025-02-26T05:04:01Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/workers	604.231s
     | time="2025-02-26T05:10:12Z" level=error msg=FAIL
FAIL | 	github.com/photoprism/photoprism/internal/workers/auto	604.270s
```


# Inconsistencies Discovered.

UpdateSubjectCovers updates 6 x 2 for SQLite and PostgreSQL.  But 6 and 0 for MariaDB.  Executing the captured SQL against MariaDB results in 6 and 6.  Unsure if this is a defect in Gorm not reporting the number of records affected correctly, or something else.  This branch has not changed the MariaDB query.  
  
