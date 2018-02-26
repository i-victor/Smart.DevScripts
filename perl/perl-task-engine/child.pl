#!/usr/bin/env perl

# [SmartFramework / Task Engine / Child]
# (c) 2006-2016 unix-world.org - all rights reserved
# r.161028

######################################## PERL MODULES

use strict;
use warnings;
use Cwd;
use Time::HiRes;
use Term::ANSIColor;

######################################## TERM COLORS

my $clr_error = ['red'];
my $clr_warn = ['bright_magenta'];
my $clr_notice = ['bright_blue'];
my $clr_ok = ['blue'];

######################################## PARSE INI SETTINGS {{{SYNC-PERL-INIPARSE}}} ### DO NOT EDIT THIS SCRIPT !!! ### USE child.ini to store all settings ###

my %inisett = ();
my $cfname;
my $cfval;
open(CONF, "child.ini") || die "Failed to open config file child.ini";
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

my $task_script = "";
if($inisett{'ChildTask'} ne "") {
	$task_script = $inisett{'ChildTask'};
}

my $task_verbose = "yes";
if($inisett{'Verbose'} eq "no") {
	$task_verbose = "no";
}

######################################## RUNTIME

# Child Usage Info
my $num_args = $#ARGV + 1;
if($num_args != 1) {
	print colored($clr_error, "CHILD.ERR: Usage: child.pl TaskID");
	print "\n";
	exit;
}

# Test Task ID
my $the_task_id = "";
my $fArg = $ARGV[0];
$fArg =~ s/^\s+|\s+$//g; #trim
if($fArg =~ /^[_A-Za-z0-9\-]*$/) { # test if in A-Z a-z 0-9 _ -
	$the_task_id = $fArg;
}
if($the_task_id eq "") {
	print colored($clr_error, "CHILD.ERR: Task ID is EMPTY or INVALID: [ID=".$fArg."]");
	print "\n";
	exit;
}

# Test if semaphore exists
my $dir = getcwd;
my $the_child_pid = $dir."/child-semaphores/semaphore@".$the_task_id.".pid";
if(-e $the_child_pid) {
	print colored($clr_notice, "CHILD.NOTICE: Semaphore already exist, skipping ... [ID=".$the_task_id."]");
	print "\n";
	exit;
}

# create the semaphore
system("echo ".$the_task_id." > ".$the_child_pid);
if(-e $the_child_pid) {
	#OK
} else {
	print colored($clr_error, "CHILD.ERR: FAILED to Lock the Semaphore ... [ID=".$the_task_id."]");
	print "\n";
	exit;
}
END {
	# remove the semaphore if exists
	if(-e $the_child_pid) {
		system("rm ".$the_child_pid);
		if(-e $the_child_pid) {
			print colored($clr_error, "CHILD.ERR: FAILED to Unlock the Semaphore ... [ID=".$the_task_id."]");
			print "\n";
		} else {
			if($task_verbose ne "no") {
				print colored($clr_ok, "CHILD.INF: Task DONE [ID=".$the_task_id."]");
				print "\n";
			}
		}
	} else {
		print colored($clr_error, "CHILD.ERR: FAILED to Find the Semaphore for Unlocking ... [ID=".$the_task_id."]");
		print "\n";
	}
}

# Run Task script (if empty run a test only)
if($task_script eq "") {
	if($task_verbose ne "no") {
		print colored($clr_warn, "CHILD.WARN: No Task Defined in child.ini (Running TEST ONLY) !!!");
		print "\n";
	}
} else {
	if($task_verbose ne "no") {
		print colored($clr_ok, "CHILD.INF: Task FORK :: Smart.Task.Engine [ID=".$the_task_id."] "."[Script=".$task_script."]");
		print "\n";
	}
	system($task_script." ".$the_task_id); # RUN REAL TASK
	Time::HiRes::sleep(0.5);
}

######################################## EXIT

exit 0;

#END
