
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


# project definition
project(
	'vala-project', 'c', 'vala',
	version : '1.0',
	license : [ 'BSD'],
	meson_version : '>= 0.50.0',
	default_options : [ 'backend=ninja', 'warning_level=1', 'buildtype=release', 'optimization=2', 'strip=true', 'cpp_std=c++11', 'prefix=/opt/unixman' ]
)

# define min vala version required
vala_version_min = '0.44'

# define vala packages deps ; # Sample (Direct Compile): valac --pkg gtk+-3.0 --pkg webkit2gtk-4.0 "%f"
valadeps = [
	dependency('gtk+-3.0'),
	dependency('webkit2gtk-4.0'),
]

# define executable ; Sample: executable('prog', 'main.vala', 'class1.vala', 'class2.vala', dependencies : valadeps)
executable('webkit-sample', 'webkit-sample.vala', dependencies: valadeps)

# check min vala version required
vala = meson.get_compiler('vala')
if not vala.version().version_compare('>= @0@'.format(vala_version_min))
	error('Vala Compiler valac >= @0@ required !'.format(vala_version_min))
endif


###
# Meson Reference Samples for Vala: https://github.com/mesonbuild/meson/tree/master/test%20cases/vala
###

# END
