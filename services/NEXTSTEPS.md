# Services Next Steps

## Completed
- ✅ FileStore service with local filesystem backend
- ✅ R2 FileStore service with presigned URLs
- ✅ Screenshot indexing pipeline with batch processing
- ✅ Optimistic locking for WorldData updates
- ✅ BackendWorldsService abstraction for GORM and FS implementations
- ✅ WorldDataUpdater interface for storage-agnostic operations
- ✅ Google Cloud Datastore backend (gaebe package) for App Engine deployment
  - Uses protoc-gen-dal generated entities with proper datastore_tags
  - Supports composite indexes for needs_indexing queries
  - Cross-entity transactions for SaveMoveGroup atomicity
  - Backend selectable at runtime via WORLDS_SERVICE_BE/GAMES_SERVICE_BE=gae
- ✅ Consistent world ID normalization (lowercase) across all backends
  - `NormalizeWorldID()` helper in services package
  - Applied to GetWorld, UpdateWorld, DeleteWorld, CreateWorld in all backends
- ✅ PropertyLoadSaver support for Datastore entities with map fields
  - Added `implement_property_loader: true` to WorldDataDatastore and GameStateDatastore protos
  - Enables serialization of map[string]* fields as JSON blobs in Cloud Datastore
- ✅ Fixed empty world IDs in Datastore listings
  - Id field has `datastore:"-"` tag, must be populated from entity key name
  - Fixed in ListWorlds, getWorldAndData, and UpdateWorld for gaebe backend
- ✅ Resend email integration for transactional emails
  - ResendEmailSender implements oneauth.SendEmail interface
  - Conditional: uses Resend when RESEND_API_KEY is set, falls back to ConsoleEmailSender
- ✅ Fixed Connect auth: `injectAuthMetadata()` was using `AppendToOutgoingContext` instead of `NewIncomingContext`, causing user IDs to never reach gRPC services through the Connect HTTP path
- ✅ Connect auth integration tests (`web/server/connect_auth_integration_test.go`)
  - Full HTTP pipeline: Bearer token -> oneauth middleware -> Connect adapter -> gRPC incoming metadata -> mock service
  - 5 test cases covering auth propagation, no-auth, invalid token, multi-user isolation, all write endpoints

## TODO

### FileStore
- [ ] Implement file size limits and content-type validation
- [ ] Add file metadata caching to avoid repeated HeadObject calls
- [ ] Consider adding file versioning support
- [ ] Add cleanup/garbage collection for orphaned screenshots

### Screenshot Indexing
- [ ] Implement worker pool for parallel screenshot generation
- [ ] Add rate limiting to prevent overwhelming filestore
- [ ] Implement retry logic for failed theme renders
- [ ] Add proactive re-indexing for items with NeedsIndexing=true (periodic checker)
- [ ] Support screenshot generation for games (currently only worlds)

### CLI
- [ ] Implement `ww worlds` commands (list, get, show) - GitHub issue #98
  - `ww worlds list` - List worlds on active profile
  - `ww worlds get <id>` - Get world details
  - `ww worlds show <id>` - Render world map as PNG
  - Should default to active profile when no profile specified

### Testing
- [ ] Add unit tests for path security (directory traversal attempts)
- [ ] Add integration tests for screenshot pipeline
- [ ] Test optimistic locking conflicts with concurrent updates
- [ ] Test R2 presigned URL generation and expiry

### Known Issues
- [ ] Games service ID population: GameDatastore.Id also has `datastore:"-"` - likely has same empty ID bug as worlds had
- [ ] ListWorlds owner_id filtering: All three backends ignore owner_id parameter

### Documentation
- [ ] Document screenshot URL structure and theme naming
- [ ] Add examples for FileStore API usage
- [ ] Document WorldDataUpdater interface contract
