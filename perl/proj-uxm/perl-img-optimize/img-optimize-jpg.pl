#!/usr/bin/env perl

# [JPEG Optimizer :: using ImageMagick]
# (c) 2006-2016 unix-world.org - all rights reserved
# r.170704


######################################## PERL MODULES

use strict;
use warnings;
use Cwd;
use Time::HiRes;
use Term::ANSIColor;

######################################## TERM COLORS

my $clr_critical_error = ['bold bright_white on_red'];
my $clr_error = ['bold bright_white on_bright_red'];
my $clr_warn = ['bold bright_white on_yellow'];
my $clr_notice = ['bold black on_bright_cyan'];
my $clr_ok = ['bold bright_white on_green'];

######################################## CHECK ARGUMENTS

my $num_args = $#ARGV + 1;
if($num_args != 1) {
	print colored($clr_critical_error, "JPEG.OPTIMIZE:ERR : FAIL.STOP: The script ID parameter is missing (script must have 1 parameter) ...");
	print "\n";
	exit;
}
my $jpeg = $ARGV[0];

######################################## RUNTIME

if($jpeg eq "") {
	print colored($clr_error, "JPEG.OPTIMIZE:ERR : FAIL.STOP: The script ID parameter is EMPTY ...");
	print "\n";
	exit;
}

$jpeg = str_single_quotes_escapeshellarg($jpeg);

my $optimize = -1;
my $params = "";
my $quality = 0;
my $sampling = '';
$quality = `identify -format '%Q' '${jpeg}'`;
$sampling = `identify -format '%[jpeg:sampling-factor]' '${jpeg}'`;


######################################## EVAL RESULT

$quality = 0 + ($quality||0); # fix for non-integer comparison crash
print colored($clr_notice, "JPEG.OPTIMIZE:NOTICE # ${jpeg} : DETECTED QUALITY IS: ${quality}");
print "\n";
print colored($clr_notice, "JPEG.OPTIMIZE:NOTICE # ${jpeg} : DETECTED SAMPLING IS: ${sampling}");
print "\n";

if($quality > 89) {
	$params = "${params} -quality 89";
}
if($sampling ne "2x2,1x1,1x1") {
	$params = "${params} -sampling-factor 4:2:0"; # 2x2,1x1,1x1
}

if($params ne "") {
	$params = "${params} -interlace JPEG -colorspace RGB";
	$optimize = `mogrify ${params} ${jpeg}`;
	$optimize = 0 + ($optimize||0);
	if($optimize == 0) {
		print colored($clr_ok, "JPEG.OPTIMIZE:OK # ${jpeg} # ${params}");
	} else {
		print colored($clr_error, "JPEG.OPTIMIZE:NOT-OK # Failed to Optimize Image: ${jpeg} # ${params}");
	}
	print "\n";
} else {
	print colored($clr_warn, "JPEG.OPTIMIZE:SKIP # ${jpeg} # Already Optimized");
	print "\n";
}

######################################## INTERNAL FUNCTIONS {{{SYNC-PERL-FXS}}}

sub str_single_quotes_escapeshellarg {
	my $arg = shift;
	$arg =~ s/'/'\\''/g; # escape single quotes
	return "".$arg;
}

######################################## EXIT

exit 0;

#END
