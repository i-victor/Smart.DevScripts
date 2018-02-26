#!/usr/bin/env perl

# [SmartFramework / Task Engine / Daemon]
# (c) 2006-2016 unix-world.org - all rights reserved
# r.161028

# Main Runtime - Daemon (Press CTRL+C to break/stop) !

######################################## PERL MODULES

use strict;
use warnings;
use Cwd;
use Term::ANSIColor;

######################################## TERM COLORS

my $clr_error = ['bold black on_bright_red'];
my $clr_info = ['bold bright_white on_black'];
my $clr_notice = ['bold white on_bright_black'];

######################################## PARSE INI SETTINGS {{{SYNC-PERL-INIPARSE}}} ### DO NOT EDIT THIS SCRIPT !!! ### USE daemon.ini to store all settings ###

my %inisett = ();
my $cfname;
my $cfval;
open(CONF, "daemon.ini") || die "Failed to open config file daemon.ini";
while(<CONF>) {
	s/\r|\n//g;
	if (/^#/ || !/\S/) {
		next;
	}
	/^([^=]+)=(.*)$/;
	$cfname = $1;
	$cfval = $2;
	$cfname =~ s/^\s+//g;
	$cfname =~ s/\s+$//g;
	$cfval =~ s/^\s+//g;
	$cfval =~ s/\s+$//g;
	$inisett{$cfname} = $cfval;
}
close(CONF);

######################################## RUNTIME

# setting: set the min page offset (loops) for using to batch count = startRows / batchSize ; if > 0 will use (and start here) the navigation offset * number of pages by appending to URL ofs=123
my $min_offset = 0;
if($inisett{'MinPagesOffset'} > 0) {
	$min_offset = $inisett{'MinPagesOffset'};
}

# setting: set the max page offset (loops) for using to batch count = totalRows / batchSize ; if > 0 will use (and stop here) the navigation offset * number of pages by appending to URL ofs=123
my $max_offset = $inisett{'MaxPagesOffset'}; # 0
if($max_offset < 0) {
	$max_offset = 0;
}
if($max_offset <= $min_offset) {
	$max_offset = 0;
	$min_offset = 0;
}

# setting: sleep time adjustment
my $sleep_time = $inisett{'SleepTime'}; # 15
if($sleep_time < 1) {
	$sleep_time = 1;
}

######################################## RUNTIME

my $navOffset = $min_offset;

# build the reference to master.pid
my $dir = getcwd;
my $the_master_pid = $dir."/master-semaphore.pid";

# Catch CTRL+C
my $shutdown = "";
$SIG{INT} = sub { $shutdown = "stop"; exit; };

# Info
print colored($clr_info, "Smart.Task.Engine // DAEMON: Starting daemon.pl wrapper for master.pl with interval [".$sleep_time." seconds]");
print "\n";

# Main Loop
while($shutdown eq "") {
	# if master.pid semaphore exists, wait to exit
	if(-e $the_master_pid) {
		print colored($clr_error, "DAEMON.ERR: the daemon.pl process is spawning the master.pl too fast, try to adjust the sleep time in daemon.pl ... skiping the run of master.pl until next cycle ...");
		print "\n";
	} else {
		if($max_offset > 0) {
			print colored($clr_notice, "DAEMON.INF: launching master.pl [using PageOffset=".$navOffset." / MinPagesOffset=".$min_offset." / MaxPagesOffset=".$max_offset."] ...");
			print "\n";
			system($dir."/master.pl ".$navOffset);
			$navOffset = $navOffset + 1;
		} else {
			print colored($clr_notice, "DAEMON.INF: launching master.pl ...");
			print "\n";
			system($dir."/master.pl");
		}
	}
	sleep($sleep_time); # sleep time in seconds
	if($max_offset > 0) {
		if($navOffset >= $max_offset) {
			print colored($clr_notice, "DAEMON.INF: Max Offset Loops Reached: ".$max_offset.' ###');
			print "\n";
			exit;
		}
	}
}

######################################## EXIT

exit 0;

#END
