## kubent Release Notes

### 0.3.2 (2020-09-22)

**Fixes**:
- Fixes missing resources with File Collector (#40)

**Internal/Misc**:
- Fixes git sha in binary (#41)
- Increased test coverage (storeCollector) and introduced K8s fake client tests (ClusterCollector) (#37)

### 0.3.1 (2020-09-04)

**Fixes**:
- Fix missing resources (#34)
- Fix panic when collector fails to initialize (#32)

### 0.3.0 (2020-08-11)

**Features**:
- Added stdin support (#19)
- Added support for reading manifests from files (#15) 
- Improved Error reporting (#20)

**Fixes**:
- Support resources without namespace (#14) 

**Internal/misc**:
- Added first and many other tests (#9, #19), thanks @david-doit-intl 🚀
- Cleaned up logic in main (#16) 
- Added release notes and minor deprecated message improvement (#22)

### 0.2.1 (2020-05-24)

**Features**:
- Added install script (#3, #4)

**Internal/misc**:
- Moved logic to unpack last-applied config to collector (7)

### 0.2.0 (2020-04-15)

**Fixes**:
- Produce static binaries (#2) 

### 0.1.0 - Initial Release (2020-04-09)
