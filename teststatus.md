# Synopsis  
GORM v2 has introduced foreign keys to the database.  
This has caused a number of the tests to no longer function.  
The fundamental issue is that these tests worked in the past as GORM v1 did not check if the parent record existed before saving the child record.  
eg.  
You could save a PhotoLabel with random numbers as the PhotoID and LabelID without an error being raised.  This is no longer possible.  

# Tests changed so that they work with GORM v2

All the TestMains that utilise the database have been changed so they have a database hosted MUTEX.  
This is to prevent 2 or more sets of database tests running at the same time.
When this happens then 2nd or subsequent test will truncate all the data as they other test(s) are running causing random test failures.  
In testing up to four separate testing threads have attempted to run against the database at the same time using the makefile.  
The issue is the requirement to clear and refresh the unit test data so each suite of tests work correctly.  

The TestMains have an additional check for error records in the Errors table, and will mark the test suite as failed if any new records are reported.  
This is to ensure that the checking that has been implemented in Photo.Save (and will be implemented in other Save's as appropriate) hasn't found any scenarios where mismatches between the in Struct ID field and the sub Struct's ID field are detected.  


## Have to create records in the database or the tested function will fail due to foreign key violations  
| File | Test |
|----------------------------------------|---------------------------------------------|
| internal/entity/photo_label_test.go | TestPhotoLabel_Save/success |
| internal/entity/photo_label_test.go | TestPhotoLabel_Save/photo not nil and label not nil |
| internal/entity/photo_album_test.go | TestFirstOrCreatePhotoAlbum/not yet existing photo_album |
| internal/entity/photo_album_test.go | TestPhotoAlbum_Save/success |
| internal/entity/file_test.go | TestFile_Create/file already exists |
| internal/entity/file_test.go | TestFile_Update/success |
| internal/entity/file_test.go | TestFile_Delete/permanently |
| internal/entity/file_test.go | TestFile_Delete/not permanently |
| internal/entity/file_sync_test.go | TestNewFileSync |
| internal/entity/file_sync_test.go | TestFirstOrCreateFileSync/not yet existing |
| internal/entity/file_sync_test.go | TestFileSync_Updates/success |
| internal/entity/file_sync_test.go | TestFileSync_Update/success |
| internal/entity/file_sync_test.go | TestFileSync_Save/success |
| internal/entity/file_share_test.go | TestFirstOrCreateFileShare/not yet existing |
| internal/entity/file_share_test.go | TestFirstOrCreateFileShare/existing |
| internal/entity/file_share_test.go | TestFileShare_Updates/success |
| internal/entity/file_share_test.go | TestFileShare_Update/success |
| internal/entity/file_share_test.go | TestFileShare_Save/success |
| internal/entity/auth_user_share_test.go | TestUserShare_Create |
| internal/entity/auth_user_settings_test.go | TestCreateUserSettings/Success |
| internal/entity/auth_user_details.go | TestCreateUserDetails/Success |



## Have to populate extra fields or the tested function will fail (or pass as the errors aren't checked in the test) due to foreign key violations  
| File | Test |
|----------------------------------------|---------------------------------------------|
| internal/entity/keyword_test.go | TestMarker_ClearFace/empty face id |
| internal/entity/file_sync_test.go | TestFirstOrCreateFileSync/existing |
| internal/entity/auth_user_details_test.go | TestUserDetails_Updates |
| internal/entity/auth_user_settings_test.go | TestUserSettings_Updates |
| internal/entity/file_test.go | TestFile_Undelete/success |
| internal/entity/file_test.go | TestFile_Undelete/file not missing |



## Have to provide an ID value or tested function will fail with Where clause missing  
| File | Test |
|----------------------------------------|---------------------------------------------|
| internal/entity/keyword_test.go | TestKeyword_Update/success |


## Specials  
| File | Test | Description |
|----------------------------------------|---------------------------------------------|------------------------------------------------------|
| internal/entity/search/photos_filter_filter_test.go | TestPhotosFilterFilter/CenterPercent | Soft delete a record that was "hidden" due to duplicate ID values in Fixture |
| internal/entity/search/photos_filter_filter_test.go | TestPhotosQueryFilter/CenterPercent | Soft delete a record that was "hidden" due to duplicate ID values in Fixture |
| internal/entity/query/moments_test.go | TestRemoveDuplicateMoments/Ok | sqlite issue on GORMv1 which hasn't shown up on GORMv2 |
| internal/entity/query/files_test.go | TestFilesByUID/error | GORMv1 vs GORMv2 differences |
| internal/entity/query/file_selection_test.go | TestFileSelection/ShareSelectionOriginals | Not sure why MediaType is empty string on GORMv1 as it shouldn't be, so force it to ensure test works.  See Filefixture which sets file name with .jpg and BeforeCreate which sets the MediaType. |
| internal/entity/entity_update_test.go | TestEntitiy_Update/Photo01 | add checking that the Camera isn't removed |



## Fixed Tests
| File | Test | Description |
|----------------------------------------|---------------------------------------------|------------------------------------------------------|
| internal/entity/auth_user_details_test.go | TestUserDetails_Updates | Validate that no errors are returned |
| internal/entity/auth_user_settings_test.go | TestUserSettings_Updates | Validate that no errors are returned |
| internal/entity/photo_test.go | TestPhoto_ClassifyLabels/NewPhoto | Use a new photo struct, and load data correctly.  (Name of this test is misleading) |
| internal/entity/photo_test.go | TestPhoto_ClassifyLabels/ExistingPhoto | Use a new photo struct, and load data correctly. (Name of this test is misleading) |



## New Tests
| File | Test | Description |
|----------------------------------------|---------------------------------------------|------------------------------------------------------|
| internal/entity/keyword_test.go | TestKeyword_UpdateNoID/success | Validates "id value required but not provided" error for Update |
| internal/entity/keyword_test.go | TestKeyword_Updates/success ID on keyword | Validate that Update saves a Keyword when the ID is in the struct |
| internal/entity/keyword_test.go | TestKeyword_Updates/failure | Validate that Update fails when an ID is not in the struct or the request |
| internal/entity/marker_test.go | TestMarker_ClearFace/missing markeruid | Validates "markeruid required but not provided" error for Update |
| internal/entity/marker_test.go | TestMarker_Matched/missing markeruid | Validates "markeruid required but not provided" error for Update |
| internal/entity/auth_user_test.go | TestUser_ValidatePreload/* | Validates that Preload is used to get child attributes |
| internal/entity/query/files_test.go | TestFilesByUID/Negative limit with offset | Validates limits and offsets |
| internal/entity/query/files_test.go | TestFilesByUID/offset and limit | Validates limits and offsets |
| internal/entity/dbtest/dbtest_init_test.go | TestInit/* | checks that the number of records in a fresh database is correct |
| internal/entity/dbtest/dbtest_foreignkey_test.go | TestDbtestForeignKey_Validate/Photos_CameraID | makes sure that the database throws a foreign key error |
| internal/entity/dbtest/dbtest_foreignkey_test.go | TestDbtestForeignKey_Validate/Photos_LensID | makes sure that the database throws a foreign key error |
| internal/entity/dbtest/dbtest_fieldlength_test.go | TestInitDBLengths/PhotoMaxVarLengths | makes sure that the database can hold specified maximum length of each column in a Photo |
| internal/entity/dbtest/dbtest_fieldlength_test.go | TestInitDBLengths/PhotoExceedMax* | makes sure that the database throws an error when the maximum length is exceeded by 1 character in a Photo |
| internal/entity/dbtest/dbtest_blocking_test.go | TestEntity_UpdateDBErrors | verifies that entity.Update detects and returns database level errors |
| internal/entity/dbtest/dbtest_blocking_test.go | TestEntity_SaveDBErrors | verifies that entity.Save detects and returns database level errors |
| internal/entity/dbtest/dbtest_migration_test.go | TestDialectSQLite3/ValidMigration | verifies that the migration has completed and that auto_increment is working, and fields can be populated with min and max values |
| internal/entity/dbtest/dbtest_migration_test.go | TestDialectSQLite3/InvalidDataUpgrade | verifies that a database with invalid foreign keys will be cleansed and then migrated |
| internal/entity/dbtest/dbtest_migration_test.go | TestDialectMysql/ValidMigration | verifies that the migration has completed and that auto_increment is working, and fields can be populated with min and max values |
| internal/entity/dbtest/dbtest_migration_test.go | TestDialectMysql/InvalidDataUpgrade | verifies that a database with invalid foreign keys will be cleansed and then migrated |
| internal/entity/dbtest/dbtest_validatecreatesave_test.go | * | A set of tests that compare Gorm version functionality.  They can show you what to expect for a variety of scenarios.  The Entity.Save test no longer fails as Entity.Save has been uplifted to work like Gorm V1. |
| internal/entity/camera_test.go | TestCamera_ScopedSearchFirst/* | Validate that ScopedSearchFirstCamera returns the expected results/errors |
| internal/entity/entitiy_save_test.go | TestSave/NewParentPhotoWithNewChildDetails | Validate that FK violations do not happen when saving Details with a new Photo |
| internal/entity/entity_update_test.go | TestEntitiy_Update/InconsistentCameraVSCameraID | Validates that an inconsistent CameraID and Camera.ID is handled |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/NoInterface | Validates that ModelValuesStructOption handles lack of an Interface |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/NewPhoto | Validates that ModelValuesStructOption handles a new and empty Struct |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/ExistingPhoto | Validates that ModelValuesStructOption handles a populated Struct and does/doesn't remove appropriate attributes from the Struct |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/NewFace | Validates that ModelValuesStructOption handles a new and empty Struct |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/ExistingFace | Validates that ModelValuesStructOption handles a populated Struct and does/doesn't remove appropriate attributes from the Struct |
| internal/entity/entity_values_test.go | TestModelValuesStructOption/AllTypes | Validates that ModelValuesStructOption handles a populated Struct and does/doesn't remove appropriate attributes from the Struct.  **This test actions all the known types in PhotoPrism.** |
| internal/entity/file_test.go | TestFile_MissingPhotoID/No PhotoID or Photo | Validate that an error is raised when attempting to create a File without a Photo |
| internal/entity/file_test.go | TestFile_MissingPhotoID/No PhotoID and Photo.ID = 0 | Validate that an error is raised when attempting to create a File with a Photo that hasn't been created |
| internal/entity/file_test.go | TestFile_MissingPhotoID/PhotoID = 0 and Photo.ID = 0 | Validate that an error is raised when attempting to create a File without a PhotoID and a Photo that hasn't been created|
| internal/entity/lens_test.go | TestLens_ScopedSearchFirst/* | Validate that ScopedSearchFirstLens returns the expected results/errors |
| internal/entity/photo_label_test.go | TestFirstOrCreatePhotoLabel/success path 1 | Validate that an existing Label is added |
| internal/entity/photo_label_test.go | TestFirstOrCreatePhotoLabel/success path 2 | Validate that a new Label is added |
| internal/entity/photo_quality_test.go | TestPhoto_QualityScore/digikam test | Test scenario where a new Photo is created with and saved correctly, then a 2nd file is added and saved.  Ensure that the QualityScore is not incorrect.  This replicates a front end acceptance test that was failing for GormV2 |
| internal/entity/photo_test.go | TestSavePhotoForm/BadCamera | Validate that when a bad CameraID is passed from a form it is replaced with UnknownCameraID |
| internal/entity/photo_test.go | TestSavePhotoForm/BadLens | Validate that when a bad LensID is passed from a form it is replaced with UnknownLensID |
| internal/entity/photo_test.go | TestPhoto_Save/BadCameraID | Validate that when a mismatch between CameraID and Camera.ID is saved, the CameraID wins and an Error is added to database |
| internal/entity/photo_test.go | TestPhoto_Save/BadCellID | Validate that when a mismatch between CellID and Cell.ID is saved, the CellID wins and an Error is added to database |
| internal/entity/photo_test.go | TestPhoto_Save/BadLensID | Validate that when a mismatch between LensID and Lens.ID is saved, the LensID wins and an Error is added to database |
| internal/entity/photo_test.go | TestPhoto_Save/BadPlaceID | Validate that when a mismatch between PlaceID and Place.ID is saved, the PlaceID wins and an Error is added to database |
| internal/entity/photo_test.go | TestPhoto_UnscopedSearch/* | Validate that UnscopedSearchPhoto returns that expected results/errors |
| internal/entity/photo_test.go | TestPhoto_ScopedSearch/* | Validate that ScopedSearchPhoto returns that expected results/errors |
| internal/entity/photos_test.go | TestPhotos_UnscopedSearch/* | Validate that UnscopedSearchPhotos returns that expected results/errors |
| internal/entity/photos_test.go | TestPhotos_ScopedSearch/* | Validate that ScopedSearchPhotos returns that expected results/errors |
| internal/photoprism/index_mediafile_test.go | TestIndex_MediaFile/twoFiles | Test scenario where 2 files are indexed (Primary and Json) that it is done correctly.  This replicates a front end acceptance test that was failing for GormV2 |
| internal/performancetest/benchmark_100k_test.go | Benchmark100k_SQLite/* | Benchmark's for SQLite with a 100k record database |
| internal/performancetest/benchmark_100k_test.go | Benchmark100k_MySQL/* | Benchmark's for Mariadb with a 100k record database |
| internal/performancetest/benchmark_migration_test.go | BenchmarkMigration_SQLite/* | Database Migration Benchmark's for SQLite with a 100k record database |
| internal/performancetest/benchmark_migration_test.go | BenchmarkMigration_MySQL/* | Database Migration Benchmark's for MySQL with a 100k record database |

**Please note that the tests in internal/entity/dbtest all MUST use the dbtestMutex as they must run synchronous due to the database blocking tests.  Failure to include the dbtestMutex will cause unexpected failure of the test.**  


# Testing Status  

## Overview  

All tests that pass against the base develop branch of PhotoPrism are passing against gorm2 branch of PhotoPrism.  
There is one acceptance test that is failing in both branches:  
 ✖ Common: Add/Remove Photos to/from album  

Benchmark tests has been created and executed against base develop and gorm2 branch of PhotoPrism.
There are not significant performance differences between the two branches, with some database functions faster and some slower.  
An example of slower is the creation and deletion of a Photo (and all it's associated child records) which on MariaDB has gone from 33ms to 36ms (ms = milli second).  

## Unit Test Details
The following is the detailed status of unit testing against sqlite.  

```
Chrome Headless 133.0.0.0 (Linux x86_64): Executed 347 of 347 SUCCESS (0.323 secs / 0.075 secs)
TOTAL: 347 SUCCESS
TOTAL: 347 SUCCESS

=============================== Coverage summary ===============================
Statements   : 69.3% ( 1991/2873 )
Branches     : 50.93% ( 1035/2032 )
Functions    : 68.09% ( 476/699 )
Lines        : 69.75% ( 1910/2738 )
================================================================================
```

Removing test database files...  
find ./internal -type f -name ".test.*" -delete  
Running all Go tests...  
richgo test -parallel 1 -count 1 -cpu 1 -tags="slow,develop" -timeout 20m ./pkg/... ./internal/...  

| Status | Path/Test |
| ------ | --------------------------------------------------------------------- |
| PASS | github.com/photoprism/photoprism/pkg/authn |
| PASS | github.com/photoprism/photoprism/pkg/capture |
| PASS | github.com/photoprism/photoprism/pkg/checksum |
| PASS | github.com/photoprism/photoprism/pkg/clean |
| PASS | github.com/photoprism/photoprism/pkg/clusters |
| PASS | github.com/photoprism/photoprism/pkg/env |
| PASS | github.com/photoprism/photoprism/pkg/fs |
| PASS | github.com/photoprism/photoprism/pkg/fs/fastwalk |
| PASS | github.com/photoprism/photoprism/pkg/geo |
| PASS | github.com/photoprism/photoprism/pkg/geo/pluscode |
| PASS | github.com/photoprism/photoprism/pkg/geo/s2 |
| PASS | github.com/photoprism/photoprism/pkg/i18n |
| PASS | github.com/photoprism/photoprism/pkg/list |
| PASS | github.com/photoprism/photoprism/pkg/log/dummy |
| PASS | github.com/photoprism/photoprism/pkg/log/level |
| PASS | github.com/photoprism/photoprism/pkg/media |
| PASS | github.com/photoprism/photoprism/pkg/media/colors |
| PASS | github.com/photoprism/photoprism/pkg/media/http/header |
| PASS | github.com/photoprism/photoprism/pkg/media/http/scheme |
| PASS | github.com/photoprism/photoprism/pkg/media/projection |
| PASS | github.com/photoprism/photoprism/pkg/media/video |
| PASS | github.com/photoprism/photoprism/pkg/react |
| PASS | github.com/photoprism/photoprism/pkg/rnd |
| PASS | github.com/photoprism/photoprism/pkg/time/unix |
| PASS | github.com/photoprism/photoprism/pkg/txt |
| PASS | github.com/photoprism/photoprism/pkg/txt/clip |
| PASS | github.com/photoprism/photoprism/pkg/txt/report |
| PASS | github.com/photoprism/photoprism/pkg/vector |
| PASS | github.com/photoprism/photoprism/internal/ai/classify |
| PASS | github.com/photoprism/photoprism/internal/ai/face |
| PASS | github.com/photoprism/photoprism/internal/ai/nsfw |
| SKIP | github.com/photoprism/photoprism/internal/entity/legacy	[no test files] |
| PASS | github.com/photoprism/photoprism/internal/api |
| PASS | github.com/photoprism/photoprism/internal/auth/acl |
| PASS | github.com/photoprism/photoprism/internal/auth/oidc |
| PASS | github.com/photoprism/photoprism/internal/auth/session |
| PASS | github.com/photoprism/photoprism/internal/commands |
| PASS | github.com/photoprism/photoprism/internal/config |
| PASS | github.com/photoprism/photoprism/internal/config/customize |
| PASS | github.com/photoprism/photoprism/internal/config/pwa |
| PASS | github.com/photoprism/photoprism/internal/config/ttl |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/apple	[no test files] |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/intel	[no test files] |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/nvidia	[no test files] |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/v4l	[no test files] |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/vaapi	[no test files] |
| PASS | github.com/photoprism/photoprism/internal/entity |
| SKIP | github.com/photoprism/photoprism/internal/testextras	[no test files] |
| PASS | github.com/photoprism/photoprism/internal/entity/dbtest |
| PASS | github.com/photoprism/photoprism/internal/entity/migrate |
| PASS | github.com/photoprism/photoprism/internal/entity/query |
| PASS | github.com/photoprism/photoprism/internal/entity/search |
| PASS | github.com/photoprism/photoprism/internal/entity/search/viewer |
| PASS | github.com/photoprism/photoprism/internal/entity/sortby |
| PASS | github.com/photoprism/photoprism/internal/event |
| PASS | github.com/photoprism/photoprism/internal/ffmpeg |
| PASS | github.com/photoprism/photoprism/internal/ffmpeg/encode |
| PASS | github.com/photoprism/photoprism/internal/form |
| PASS | github.com/photoprism/photoprism/internal/functions |
| PASS | github.com/photoprism/photoprism/internal/meta |
| PASS | github.com/photoprism/photoprism/internal/mutex |
| PASS | github.com/photoprism/photoprism/internal/performancetest |
| PASS | github.com/photoprism/photoprism/internal/photoprism |
| PASS | github.com/photoprism/photoprism/internal/photoprism/backup |
| PASS | github.com/photoprism/photoprism/internal/photoprism/get |
| PASS | github.com/photoprism/photoprism/internal/server |
| PASS | github.com/photoprism/photoprism/internal/server/limiter |
| PASS | github.com/photoprism/photoprism/internal/server/process |
| PASS | github.com/photoprism/photoprism/internal/server/wellknown |
| PASS | github.com/photoprism/photoprism/internal/service |
| PASS | github.com/photoprism/photoprism/internal/service/hub |
| PASS | github.com/photoprism/photoprism/internal/service/hub/places |
| PASS | github.com/photoprism/photoprism/internal/service/maps |
| PASS | github.com/photoprism/photoprism/internal/service/webdav |
| PASS | github.com/photoprism/photoprism/internal/thumb |
| PASS | github.com/photoprism/photoprism/internal/thumb/avatar |
| PASS | github.com/photoprism/photoprism/internal/thumb/crop |
| PASS | github.com/photoprism/photoprism/internal/thumb/frame |
| PASS | github.com/photoprism/photoprism/internal/workers |
| PASS | github.com/photoprism/photoprism/internal/workers/auto |



The following is the status of unit testing against mariadb, which drops the database as part of the init.  
Resetting acceptance database...  
mysql < scripts/sql/reset-acceptance.sql  
Running all Go tests on MariaDB...  
PHOTOPRISM_TEST_DRIVER="mysql" PHOTOPRISM_TEST_DSN="root:photoprism@tcp(mariadb:4001)/acceptance?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true" richgo test -parallel 1 -count 1 -cpu 1 -tags="slow,develop" -timeout 20m ./pkg/... ./internal/...  
| Status | Path/Test |
| ------ | --------------------------------------------------------------------- |
| PASS | github.com/photoprism/photoprism/pkg/authn |
| PASS | github.com/photoprism/photoprism/pkg/capture |
| PASS | github.com/photoprism/photoprism/pkg/checksum |
| PASS | github.com/photoprism/photoprism/pkg/clean |
| PASS | github.com/photoprism/photoprism/pkg/clusters |
| PASS | github.com/photoprism/photoprism/pkg/env |
| PASS | github.com/photoprism/photoprism/pkg/fs |
| PASS | github.com/photoprism/photoprism/pkg/fs/fastwalk |
| PASS | github.com/photoprism/photoprism/pkg/geo |
| PASS | github.com/photoprism/photoprism/pkg/geo/pluscode |
| PASS | github.com/photoprism/photoprism/pkg/geo/s2 |
| PASS | github.com/photoprism/photoprism/pkg/i18n |
| PASS | github.com/photoprism/photoprism/pkg/list |
| PASS | github.com/photoprism/photoprism/pkg/log/dummy |
| PASS | github.com/photoprism/photoprism/pkg/log/level |
| PASS | github.com/photoprism/photoprism/pkg/media |
| PASS | github.com/photoprism/photoprism/pkg/media/colors |
| PASS | github.com/photoprism/photoprism/pkg/media/http/header |
| PASS | github.com/photoprism/photoprism/pkg/media/http/scheme |
| PASS | github.com/photoprism/photoprism/pkg/media/projection |
| PASS | github.com/photoprism/photoprism/pkg/media/video |
| PASS | github.com/photoprism/photoprism/pkg/react |
| PASS | github.com/photoprism/photoprism/pkg/rnd |
| PASS | github.com/photoprism/photoprism/pkg/time/unix |
| PASS | github.com/photoprism/photoprism/pkg/txt |
| PASS | github.com/photoprism/photoprism/pkg/txt/clip |
| PASS | github.com/photoprism/photoprism/pkg/txt/report |
| PASS | github.com/photoprism/photoprism/pkg/vector |
| PASS | github.com/photoprism/photoprism/internal/ai/classify |
| PASS | github.com/photoprism/photoprism/internal/ai/face |
| PASS | github.com/photoprism/photoprism/internal/ai/nsfw |
| PASS | github.com/photoprism/photoprism/internal/api |
| PASS | github.com/photoprism/photoprism/internal/auth/acl |
| PASS | github.com/photoprism/photoprism/internal/auth/oidc |
| PASS | github.com/photoprism/photoprism/internal/auth/session |
| PASS | github.com/photoprism/photoprism/internal/commands |
| PASS | github.com/photoprism/photoprism/internal/config |
| PASS | github.com/photoprism/photoprism/internal/config/customize |
| PASS | github.com/photoprism/photoprism/internal/config/pwa |
| PASS | github.com/photoprism/photoprism/internal/config/ttl |
| PASS | github.com/photoprism/photoprism/internal/entity |
| PASS | github.com/photoprism/photoprism/internal/entity/dbtest |
| SKIP | github.com/photoprism/photoprism/internal/entity/legacy	[no test files]|
| PASS | github.com/photoprism/photoprism/internal/entity/migrate |
| PASS | github.com/photoprism/photoprism/internal/entity/query |
| PASS | github.com/photoprism/photoprism/internal/entity/search |
| PASS | github.com/photoprism/photoprism/internal/entity/search/viewer |
| PASS | github.com/photoprism/photoprism/internal/entity/sortby |
| PASS | github.com/photoprism/photoprism/internal/event |
| PASS | github.com/photoprism/photoprism/internal/ffmpeg |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/apple	[no test files]|
| PASS | github.com/photoprism/photoprism/internal/ffmpeg/encode |
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/intel	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/nvidia	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/v4l	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/vaapi	[no test files]|
| PASS | github.com/photoprism/photoprism/internal/form |
| PASS | github.com/photoprism/photoprism/internal/functions |
| PASS | github.com/photoprism/photoprism/internal/meta |
| PASS | github.com/photoprism/photoprism/internal/mutex |
| PASS | github.com/photoprism/photoprism/internal/performancetest |
| PASS | github.com/photoprism/photoprism/internal/photoprism |
| PASS | github.com/photoprism/photoprism/internal/photoprism/backup |
| PASS | github.com/photoprism/photoprism/internal/photoprism/get |
| PASS | github.com/photoprism/photoprism/internal/server |
| PASS | github.com/photoprism/photoprism/internal/server/limiter |
| PASS | github.com/photoprism/photoprism/internal/server/process |
| PASS | github.com/photoprism/photoprism/internal/server/wellknown |
| PASS | github.com/photoprism/photoprism/internal/service |
| PASS | github.com/photoprism/photoprism/internal/service/hub |
| PASS | github.com/photoprism/photoprism/internal/service/hub/places |
| PASS | github.com/photoprism/photoprism/internal/service/maps |
| PASS | github.com/photoprism/photoprism/internal/service/webdav |
| SKIP | github.com/photoprism/photoprism/internal/testextras	[no test files]|
| PASS | github.com/photoprism/photoprism/internal/thumb |
| PASS | github.com/photoprism/photoprism/internal/thumb/avatar |
| PASS | github.com/photoprism/photoprism/internal/thumb/crop |
| PASS | github.com/photoprism/photoprism/internal/thumb/frame |
| PASS | github.com/photoprism/photoprism/internal/workers |
| PASS | github.com/photoprism/photoprism/internal/workers/auto |



## Acceptance Test Details
The following is current state of acceptance testing against sqlite:  
Test Common: Disable upload, download, edit and share has been skipped as it is not currently working (as at 2025-05-25) due to "changes to some classes so that the share and webdav-share action have the same"  
This change to frontend/tests/acceptance/acceptance-auth/settings/general.js (adding .skip) has not been checked in.  


[ -f "./storage/acceptance/index.db" ] || (cd storage && rm -rf acceptance && wget -c https://dl.photoprism.app/qa/acceptance.tar.gz -O - | tar -xz)  
cp -f storage/acceptance/backup.db storage/acceptance/index.db  
cp -f storage/acceptance/config-sqlite/settingsBackup.yml storage/acceptance/config-sqlite/settings.yml  
./photoprism --auth-mode="password" -c "./storage/acceptance/config-sqlite" --test start -d  
sleep 20  
Running JS acceptance-auth tests in Chrome...  
(cd frontend &&	npm run testcafe -- "chrome --headless=new" --test-grep "^(Multi-Window)\:*" --test-meta mode=auth --config-file ./testcaferc.json --experimental-multiple-windows "tests/acceptance" && npm run testcafe -- "chrome --headless=new" --test-grep "^(Common|Core)\:*" --test-meta mode=auth --config-file ./testcaferc.json "tests/acceptance")  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --test-grep ^(Multi-Window)\:* --test-meta mode=auth --config-file ./testcaferc.json --experimental-multiple-windows tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Test link sharing  
 ✓ Multi-Window: Verify visitor role has limited permissions  
  
  
 1 passed (3m 45s)  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --test-grep ^(Common|Core)\:* --test-meta mode=auth --config-file ./testcaferc.json tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Test authentication  
 ✓ Common: Login and Logout  
 ✓ Common: Login with wrong credentials  
 ✓ Common: Change password  
 ✓ Common: Delete Clipboard on logout  
  
 Test components  
 ✓ Common: Mobile Toolbar  
  
 Test account settings  
 ✓ Core: Sign in with recovery code  
 ✓ Core: Create App Password  
 ✓ Core: Check App Password has limited permissions and last updated is set  
 ✓ Core: Try to login with invalid credentials/insufficient scope  
 ✓ Core: Delete App Password  
 ✓ Common: Try to activate 2FA with wrong password/passcode  
  
 Test general settings  
 ✓ Common: Disable delete  
 ✓ Common: Change language  
 ✓ Common: Disable pages: import, originals, logs, moments, places, library  
 ✓ Common: Disable people and labels  
 ✓ Common: Disable private, archive and quality filter  
 - Common: Disable upload, download, edit and share
  
 Test link sharing  
 ✓ Common: Create, view, delete shared albums  
  
  
 17 passed (28m 37s)  
 1 skipped  
  
 Warnings (4):  
```
 --  
  TestCafe cannot interact with the <i class="mdi-power mdi v-icon notranslate v-theme--default v-icon--size-default" aria-hidden="true"></i> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-overlay__scrim"></div>.  
  Review your code to prevent this behavior.  
  
     38 |      .click(Selector(".action-confirm"));  
     39 |  }  
     40 |  
     41 |  async logout() {  
     42 |    await menu.openNav();  
   > 43 |    await t.click(Selector("button i.mdi-power"));  
     44 |  }  
     45 |  
     46 |  async testCreateEditDeleteSharingLink(type) {  
     47 |    await menu.openPage(type);  
     48 |    const FirstAlbum = await album.getNthAlbumUid("all", 0);  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/page.js:43:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/page.js:7:42)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/page.js:7:42)  
 --  
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default text-on-surface-variant v-btn--density-default v-btn--size-default v-btn--variant-text v-tab" disabled="" tabindex="-1" role="tab" aria-selected="false" id="tab-labels" value="labels">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-slide-group__content">...</div>.  
  Review your code to prevent this behavior.  
  
     254 |    await t.click(settings.peopleCheckbox).click(settings.labelsCheckbox);  
     255 |    await t.eval(() => location.reload());  
     256 |    await menu.openPage("browse");  
     257 |    await toolbar.setFilter("view", "Cards");  
     258 |    await t.click(page.cardTitle.nth(0));  
   > 259 |    await t.click(photoedit.labelsTab);  
     260 |  
     261 |    await t.expect(photoedit.addLabel.exists).notOk();  
     262 |  
     263 |    await t.click(photoedit.peopleTab);  
     264 |  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:259:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:12:48)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:12:48)  
 --  
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default text-on-surface-variant v-btn--density-default v-btn--size-default v-btn--variant-text v-tab" disabled="" tabindex="-1" role="tab" aria-selected="false" id="tab-people" value="people">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-slide-group__content">...</div>.  
  Review your code to prevent this behavior.  
  
     258 |    await t.click(page.cardTitle.nth(0));  
     259 |    await t.click(photoedit.labelsTab);  
     260 |  
     261 |    await t.expect(photoedit.addLabel.exists).notOk();  
     262 |  
   > 263 |    await t.click(photoedit.peopleTab);  
     264 |  
     265 |    await t.expect(Selector("div.p-faces ").exists).notOk();  
     266 |  
     267 |    await t.click(photoedit.dialogClose);  
     268 |  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:263:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:12:48)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-auth/settings/general.js:12:48)  
 --  
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default bg-grey-darken-2 v-btn--density-comfortable v-btn--size-small v-btn--variant-flat action-clear" style="">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-overlay__scrim"></div>.  
  Review your code to prevent this behavior.  
  
     48 |    }  
     49 |  }  
     50 |  
     51 |  async clearSelection() {  
     52 |    await this.openContextMenu();  
   > 53 |    await t.click(Selector(".action-clear"));  
     54 |  }  
     55 |}  
     56 |  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:53:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:1:40)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:1:40)  
```
./photoprism --auth-mode="password" -c "./storage/acceptance/config-sqlite" --test stop  
cp -f storage/acceptance/backup.db storage/acceptance/index.db  
cp -f storage/acceptance/config-sqlite/settingsBackup.yml storage/acceptance/config-sqlite/settings.yml  
rm -rf storage/acceptance/sidecar/2020  
rm -rf storage/acceptance/sidecar/2011  
rm -rf storage/acceptance/originals/2010  
rm -rf storage/acceptance/originals/2020  
rm -rf storage/acceptance/originals/2011  
rm -rf storage/acceptance/originals/2013  
rm -rf storage/acceptance/originals/2017  
./photoprism --auth-mode="public" -c "./storage/acceptance/config-sqlite" --test start -d  
sleep 20  
Running public-mode tests in Chrome...  
(cd frontend &&	find ./tests/acceptance -type f -name "*.js" | xargs -i perl -0777 -ne 'while(/(?:mode: \"auth[^,]*\,)|(Multi-Window\:[A-Za-z 0-9\-_]*)/g){print "$1\n" if ($1);}' {} | xargs -I testname bash -c 'npm run testcafe -- "chrome --headless=new" --experimental-multiple-windows --test-meta mode=public --config-file ./testcaferc.json --test "testname" "tests/acceptance"'  && npm run testcafe -- "chrome --headless=new" --test-grep "^(Common|Core)\:*" --test-meta mode=public --config-file ./testcaferc.json "tests/acceptance")  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --experimental-multiple-windows --test-meta mode=public --config-file ./testcaferc.json --test Multi-Window: Test places tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Search and open photo from places  
 ✓ Multi-Window: Test places  
  
  
 1 passed (1m 01s)  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --experimental-multiple-windows --test-meta mode=public --config-file ./testcaferc.json --test Multi-Window: Navigate from card view to place tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Test photos  
 ✓ Multi-Window: Navigate from card view to place  
  
  
 1 passed (7s)  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --experimental-multiple-windows --test-meta mode=public --config-file ./testcaferc.json --test Multi-Window: Navigate from card view to photos taken at the same date tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Test photos  
 ✓ Multi-Window: Navigate from card view to photos taken at the same date  
  
  
 1 passed (13s)  
  
> photoprism@1 testcafe  
> testcafe chrome --headless=new --test-grep ^(Common|Core)\:* --test-meta mode=public --config-file ./testcaferc.json tests/acceptance  
  
 Running tests in:  
 - Chrome 133.0.0.0 / Ubuntu 24.10  
  
 Test albums  
 ✓ Common: Create/delete album on /albums  
 ✓ Common: Create/delete album during add to album  
 ✓ Common: Update album details  
 ✓ Common: Add/Remove Photos to/from album  
 ✓ Common: Use album search and filters  
 ✓ Common: Test album autocomplete  
 ✓ Common: Create, Edit, delete sharing link  
 ✓ Common: Verify album sort options  
  
 Test calendar  
 ✓ Common: View calendar  
 ✓ Common: Update calendar details  
 ✓ Common: Create, Edit, delete sharing link for calendar  
 ✓ Common: Create/delete album-clone from calendar  
 ✓ Common: Verify calendar sort options  
  
 Test components  
 ✓ Common: Test filter options  
 ✓ Common: Fullscreen mode  
 ✓ Common: Mosaic view  
 ✓ Common: List view  
 ✓ Common: Card view  
 ✓ Common: Mobile Toolbar  
  
 Test folders  
 ✓ Common: View folders  
 ✓ Common: Update folder details  
 ✓ Common: Create, Edit, delete sharing link  
 ✓ Common: Create/delete album-clone from folder  
 ✓ Common: Verify folder sort options  
  
 Test labels  
 ✓ Common: Remove/Activate Add/Delete Label from photo  
 ✓ Common: Toggle between important and all labels  
 ✓ Common: Rename Label  
 ✓ Common: Add label to album  
 ✓ Common: Delete label  
  
 Import file from folder  
 ✓ Common: Import files from folder using copy  
  
 Test index  
 ✓ Common: Index files from folder  
  
 Test moments  
 ✓ Common: Update moment details  
 ✓ Common: Create, Edit, delete sharing link for moment  
 ✓ Common: Create/delete album-clone from moment  
 ✓ Common: Verify moment sort options  
  
 Test files  
 ✓ Common: Navigate in originals  
 ✓ Common: Add original files to album  
 ✓ Common: Download available in originals  
  
 Test people  
 ✓ Common: Add name to new face and rename subject  
 ✓ Common: Add + Reject name on people tab  
 ✓ Common: Test mark subject as favorite  
 ✓ Common: Test new face autocomplete  
 ✓ Common: Remove face  
 ✓ Common: Hide face  
 ✓ Common: Hide person  
  
 Test photos archive and private functionalities  
 ✓ Common: Private/unprivate photo/video using clipboard  
 ✓ Common: Archive/restore video, photos, private photos and review photos using clipboard  
 ✓ Common: Check that archived files are not shown in monochrome/panoramas/stacks/scans/review/albums/favorites/private/videos/calendar/moments/states/labels/folders/originals  
 ✓ Common: Check that private files are not shown in monochrome/panoramas/stacks/scans/review/albums/favorites/archive/videos/calendar/moments/states/labels/folders/originals  
 ✓ Common: Check delete all dialog  
  
 Does not work in container and we have no content-disposition header anymore  
 - Common: Test download jpg file from context menu and fullscreen  
 - Common: Test download video from context menu  
 - Common: Test download multiple jpg files from context menu  
 - Common: Test raw file from context menu and fullscreen mode  
  
 Test photos upload and delete  
 ✓ Core: Upload + Delete jpg/json  
 ✓ Core: Upload + Delete video  
 ✓ Core: Upload to existing Album + Delete  
 ✓ Core: Upload jpg to new Album + Delete  
 ✓ Core: Try uploading nsfw file  
 ✓ Core: Try uploading txt file  
  
 Test photos  
 ✓ Common: Scroll to top  
 ✓ Common: Download single photo/video using clipboard and fullscreen mode  
 ✓ Common: Approve photo using approve and by adding location  
 ✓ Common: Like/dislike photo/video  
 ✓ Common: Edit photo/video  
 ✓ Common: Mark photos/videos as panorama/scan  
  
 Test about  
 ✓ Core: About page is displayed with all links  
 ✓ Core: License page is displayed with all links  
  
 Test stacks  
desktop  
 ✓ Common: View all files of a stack  
 ✓ Common: Change primary file  
 ✓ Common: Ungroup files  
 ✓ Common: Delete non primary file  
  
 Test states  
 ✓ Common: Update state details  
 ✓ Common: Create, Edit, delete sharing link for state  
 ✓ Common: Create/delete album-clone from state  
  
  
 71 passed (1h 12m 10s)  
 4 skipped  
  
 Warnings (3):  
 --  
```
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default v-btn--density-default v-btn--size-default v-btn--variant-text action-close">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-overlay__scrim"></div>.  
  Review your code to prevent this behavior.  
  
     124 |    await t.expect(photoedit.inputName.nth(0).value).eql("");  
     125 |  
     126 |    await t  
     127 |      .typeText(photoedit.inputName.nth(0), "Nicole", { replace: true })  
     128 |      .pressKey("enter")  
   > 129 |      .click(photoedit.dialogClose);  
     130 |    await contextmenu.clearSelection();  
     131 |    await t.eval(() => location.reload());  
     132 |    await t.wait(5000);  
     133 |    const PhotosInAndreaAfterRejectCount = await photo.getPhotoCount("all");  
     134 |    const Diff = PhotosInAndreaCount - PhotosInAndreaAfterRejectCount;  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-public/people.js:129:8)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-public/people.js:8:50)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/acceptance-public/people.js:8:50)  
 --  
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default bg-grey-darken-2 v-btn--density-comfortable v-btn--size-small v-btn--variant-elevated action-clear" style="">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="v-overlay__scrim"></div>.  
  Review your code to prevent this behavior.  
  
     48 |    }  
     49 |  }  
     50 |  
     51 |  async clearSelection() {  
     52 |    await this.openContextMenu();  
   > 53 |    await t.click(Selector(".action-clear"));  
     54 |  }  
     55 |}  
     56 |  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:53:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:1:40)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/context-menu.js:1:40)  
 --  
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default bg-highlight v-btn--density-default v-btn--size-default v-btn--variant-flat action-apply action-approve" disabled="">...</button> element because another element obstructs it.  
  When something overlaps the action target, TestCafe performs the action with the topmost element at the original target's location.  
  The following element with a greater z-order replaced the original action target: <div class="action-buttons">...</div>.  
  Review your code to prevent this behavior.  
  
     248 |      .typeText(this.keywords, keywords)  
     249 |      .typeText(this.notes, notes, { replace: true })  
     250 |  
     251 |      .click(Selector("button.action-approve"));  
     252 |    await t.expect(this.latitude.visible, { timeout: 5000 }).ok();  
   > 253 |    await t.click(Selector("button.action-apply")).click(Selector("button.action-close"));  
     254 |  }  
     255 |  
     256 |  async undoPhotoEdit(  
     257 |    title,  
     258 |    timezone,  
  
     at <anonymous> (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/photo-edit.js:253:13)  
     at asyncGeneratorStep (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/photo-edit.js:1:40)  
     at _next (/go/src/github.com/photoprism/photoprism/frontend/tests/acceptance/page-model/photo-edit.js:1:40)
```

## Benchmark Test Details
All the times shown in the benchmark comparisons below are in milliseconds.  

This is the tool used to compare the output from the benchmark runs.  
Please note that you may have to edit the output from the benchmark runs to remove the leading spaces and | symbol.   
```
go get golang.org/x/perf/cmd/benchstat
go run golang.org/x/perf/cmd/benchstat storage/performance_benchmark100k.gorm1.txt storage/performance_benchmark100k.txt
```

### Performance against 100k record database
As shown below, the performance between Gorm v1 and Gorm v2 is very similar, with Gorm v2 being 8.29% better overall.  
```
cd internal/performancetest
go test -skip Test -parallel 1 -count 10 -cpu 4 -failfast -tags slow -timeout 3h -benchtime 10s -bench=Benchmark100k > ../../storage/performance_benchmark100k.txt
```
goos: linux  
goarch: amd64  
pkg: github.com/photoprism/photoprism/internal/performancetest  
cpu: AMD Ryzen 7 5700X 8-Core Processor               
| Test name | storage/performance_benchmark100k.gorm1.txt | storage/performance_benchmark100k.txt | Compared |
| ---------------------------|---------------------------------|--------------------------------------- | ------------------------ |
100k_SQLite/CreateDeleteAlbum-4 | 5.411m ±  1% | 5.377m ± 2% | ~ (p=0.315 n=10)
100k_SQLite/ListAlbums-4 | 315.4m ±  1% | 259.3m ± 0% | -17.77% (p=0.000 n=10)
100k_SQLite/CreateDeleteCamera-4 | 3.326m ±  1% | 3.395m ± 1% | +2.08% (p=0.000 n=10)
100k_SQLite/CreateDeleteCellAndPlace-4 | 7.066m ±  0% | 7.046m ± 2% | ~ (p=0.393 n=10)
100k_SQLite/FileRegenerateIndex-4 | 4.576m ±  1% | 4.633m ± 1% | +1.26% (p=0.004 n=10)
100k_SQLite/CreateDeletePhoto-4 | 55.55m ±  3% | 54.63m ± 1% | -1.66% (p=0.035 n=10)
100k_SQLite/ListPhotos-4 | 341.9m ±  0% | 256.9m ± 1% | -24.86% (p=0.000 n=10)
100k_MySQL/CreateDeleteAlbum-4 | 3.651m ± 30% | 2.240m ± 3% | -38.65% (p=0.000 n=10)
100k_MySQL/ListAlbums-4 | 116.8m ±  2% | 109.0m ± 0% | -6.75% (p=0.000 n=10)
100k_MySQL/CreateDeleteCamera-4 | 1.313m ±  7% | 1.429m ± 1% | +8.79% (p=0.000 n=10)
100k_MySQL/CreateDeleteCellAndPlace-4 | 4.084m ±  1% | 3.524m ± 2% | -13.72% (p=0.000 n=10)
100k_MySQL/FileRegenerateIndex-4 | 1.210m ±  1% | 1.173m ± 4% | -3.10% (p=0.005 n=10)
100k_MySQL/CreateDeletePhoto-4 | 29.97m ±  1% | 27.39m ± 7% | -8.59% (p=0.000 n=10)
100k_MySQL/ListPhotos-4 | 451.2m ±  0% | 453.4m ± 2% | +0.49% (p=0.002 n=10)
geomean | 16.69m | 15.31m | -8.29%

### Migration performance against 100k database
The migration comparison between Gorm V1 and Gorm V2 is like comparing apples and oranges.  
The Gorm V1 migration was applying a set of standard table changes.  
The Gorm V2 migration was applying a set of standard table changes, and enabling foreign keys.  
For SQLite it was more complex as the mapped data types did not change, but the alias' data types did.  This caused Gorm to migrate each field changed one at a time (Auto), or the speedup for this (Custom).  
Please note that it is not expected for the 2nd migration via Gorm v2 to perform differently to a Gorm v1 migration, but, this has not been tested.  
```
cd internal/performancetest
go test -skip Test -parallel 1 -count 10 -cpu 4 -failfast -tags slow -timeout 3h -benchtime 1x -bench=BenchmarkMigration > ../../storage/performance_benchmarkmigration.txt
```
goos: linux  
goarch: amd64  
pkg: github.com/photoprism/photoprism/internal/performancetest  
cpu: AMD Ryzen 7 5700X 8-Core Processor               
| Test name | storage/performance_benchmarkmigration.gorm1.txt | storage/performance_benchmarkmigration.txt | Compared |
| ---------------------------|---------------------------------|--------------------------------------- | ------------------------ |
Migration_SQLite/OneKUpgradeTest_Custom-4 | 491.5m ±  9% | 1296.2m ±  7% | +163.74% (p=0.000 n=10)
Migration_SQLite/OneKUpgradeTest_Auto-4 | 484.8m ±  5% | 5176.1m ± 19% | +967.78% (p=0.000 n=10)
Migration_SQLite/TenKUpgradeTest_Custom-4 | 474.1m ±  8% | 3827.8m ± 11% | +707.46% (p=0.000 n=10)
Migration_SQLite/TenKUpgradeTest_Auto-4 | 437.5m ± 11% | 10271.8m ± 15% | +2247.69% (p=0.000 n=10)
Migration_SQLite/OneHundredKUpgradeTest_Custom-4 | 2.759 ±  3% | 18.303 ±  3% | +563.37% (p=0.000 n=10)
Migration_SQLite/OneHundredKUpgradeTest_Auto-4 | 2.717 ±  2% | 66.212 ±  2% | +2336.89% (p=0.000 n=10)
Migration_MySQL/OneKUpgradeTest-4 | 793.9m ±  6% | 7670.0m ±  9% | +866.09% (p=0.000 n=10)
Migration_MySQL/TenKUpgradeTest-4 | 836.1m ±  4% | 17672.6m ±  2% | +2013.77% (p=0.000 n=10)
Migration_MySQL/OneHundredKUpgradeTest-4 | 1.919 ±  8% | 163.363 ±  7% | +8414.10% (p=0.000 n=10)
geomean | 919.9m | 12.43 | +1251.11%
