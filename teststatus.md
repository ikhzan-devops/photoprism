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
| SKIP | github.com/photoprism/photoprism/internal/entity/legacy	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/apple	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/intel	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/nvidia	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/v4l	[no test files]|
| SKIP | github.com/photoprism/photoprism/internal/ffmpeg/vaapi	[no test files]|
| PASS | github.com/photoprism/photoprism/internal/config |
| PASS | github.com/photoprism/photoprism/internal/config/customize |
| PASS | github.com/photoprism/photoprism/internal/config/pwa |
| PASS | github.com/photoprism/photoprism/internal/config/ttl |
| PASS | github.com/photoprism/photoprism/internal/entity |
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
| SKIP | github.com/photoprism/photoprism/internal/testextras	[no test files]|
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


## Acceptance Test Details
The following is current state of acceptance testing against sqlite:  
Running auth-mode tests in Chrome...
```
> photoprism@1 testcafe
> testcafe chrome --headless=new --test-grep ^(Multi-Window)\:* --test-meta mode=auth --config-file ./testcaferc.json --experimental-multiple-windows tests/acceptance
```
Running tests in:
 - Chrome 133.0.0.0 / Ubuntu 24.10

 Test link sharing
 ✓ Multi-Window: Verify visitor role has limited permissions


 1 passed (3m 45s)

 ```
> photoprism@1 testcafe
> testcafe chrome --headless=new --test-grep ^(Common|Core)\:* --test-meta mode=auth --config-file ./testcaferc.json tests/acceptance
 ```
 
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
 ✓ Common: Disable upload, download, edit and share

 Test link sharing
 ✓ Common: Create, view, delete shared albums


 18 passed (33m 21s)

 Warnings (4):
 --
  TestCafe cannot interact with the <i class="mdi-power mdi v-icon notranslate v-theme--default v-icon--size-default" aria-hidden="true"></i> element because another element obstructs it.
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default text-on-surface-variant v-btn--density-default v-btn--size-default v-btn--variant-text v-tab" disabled="" tabindex="-1" role="tab" aria-selected="false" id="tab-labels" value="labels">...</button> element because another element obstructs it.
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default text-on-surface-variant v-btn--density-default v-btn--size-default v-btn--variant-text v-tab" disabled="" tabindex="-1" role="tab" aria-selected="false" id="tab-people" value="people">...</button> element because another element obstructs it.
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default bg-grey-darken-2 v-btn--density-comfortable v-btn--size-small v-btn--variant-flat action-clear" style="">...</button> element because another element obstructs it.


Running public-mode tests in Chrome...
```
> photoprism@1 testcafe
> testcafe chrome --headless=new --test-grep ^(Multi-Window)\:* --test-meta mode=public --config-file ./testcaferc.json --experimental-multiple-windows tests/acceptance
```
 Running tests in:
 - Chrome 133.0.0.0 / Ubuntu 24.10

 Test photos
 ✓ Multi-Window: Navigate from card view to place
 ✓ Multi-Window: Navigate from card view to photos taken at the same date

 Search and open photo from places
 ✓ Multi-Window: Test places


 3 passed (1m 15s)

```
> photoprism@1 testcafe
> testcafe chrome --headless=new --test-grep ^(Common|Core)\:* --test-meta mode=public --config-file ./testcaferc.json tests/acceptance
```
 Running tests in:
 - Chrome 133.0.0.0 / Ubuntu 24.10

 Test albums
 ✓ Common: Create/delete album on /albums
 ✓ Common: Create/delete album during add to album
 ✓ Common: Update album details
 ✖ Common: Add/Remove Photos to/from album
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


 1/71 failed (1h 10m 40s)
 4 skipped

 Warnings (3):
 --
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default v-btn--density-default v-btn--size-default v-btn--variant-text action-close">...</button> element because another element obstructs it.
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--flat v-btn--icon v-theme--default bg-grey-darken-2 v-btn--density-comfortable v-btn--size-small v-btn--variant-elevated action-clear" style="">...</button> element because another element obstructs it.
  TestCafe cannot interact with the <button type="button" class="v-btn v-btn--disabled v-btn--flat v-theme--default bg-highlight v-btn--density-default v-btn--size-default v-btn--variant-flat action-apply action-approve" disabled="">...</button> element because another element obstructs it.

## Benchmark Test Details

goos: linux
goarch: amd64
pkg: github.com/photoprism/photoprism/internal/entity
cpu: AMD Ryzen 7 5700X 8-Core Processor  

| Test name | storage/sqlite-benchmark10x.log | storage/sqlite-benchmark10x.gorm2.log | Compared |
| ---------------------------|---------------------------------|--------------------------------------- | ------------------------ |
| CreateDeleteAlbum-4 | 205.4µ ± 3%| 159.3µ ± 13% | -22.40% (p=0.000 n=10) | 
| CreateDeleteCamera-4 | 52.45µ ± 1% | 75.90µ ± 5% | +44.69% (p=0.000 n=10) | 
| CreateDeleteCellAndPlace-4 | 237.1µ ± 3% | 174.1µ ± 6% | -26.57% (p=0.000 n=10) | 
| CreateDeletePhoto-4 | 2.410m ± 6% | 2.109m ± 6% | -12.49% (p=0.000 n=10) | 
| geomean | 280.1µ | 258.1µ | -7.84% | 

goos: linux
goarch: amd64
pkg: github.com/photoprism/photoprism/internal/entity
cpu: AMD Ryzen 7 5700X 8-Core Processor 
| Test name | storage/mariadb-benchmark10x.log | storage/mariadb-benchmark10x.gorm2.log | Compared |
| ---------------------------|---------------------------------|--------------------------------------- | ------------------------ |
| CreateDeleteAlbum-4 | 2.892m ± 6% | 2.707m ± 3% | -6.39% (p=0.000 n=10) |
| CreateDeleteCamera-4 | 1.470m ± 7% | 1.536m ± 3% | +4.48% (p=0.009 n=10) |
| CreateDeleteCellAndPlace-4 | 4.564m ± 1% | 3.913m ± 9% | -14.27% (p=0.000 n=10) |
| CreateDeletePhoto-4 | 33.51m ± 8% | 36.16m ± 8% | +7.92% (p=0.007 n=10) |
| geomean 5.050m | 4.925m | -2.46% |
