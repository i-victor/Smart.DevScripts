
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


project('vala-project', 'c', 'vala')

valadeps = [
	dependency('gtk+-3.0'),
]

executable('sample', 'sample.vala', dependencies: valadeps)

# END
