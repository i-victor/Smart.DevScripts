
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


project('vala-project', 'c', 'vala')

valadeps = [
	dependency('gtk+-3.0'),
	dependency('gtksourceview-3.0'),
#	dependency('gtksourceview-4'),
]

executable('sample-editor', 'sample-editor.vala', dependencies: valadeps)

# END
