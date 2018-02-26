#!/usr/bin/env perl

# [SmartFramework / Task Engine / Master]
# (c) 2006-2016 unix-world.org - all rights reserved
# r.161028

# Fetch the Batch URL (will spawn the number of childs in $max_threads, fetching the Distribution URL: $url_get)
# This can be used in the run.pl (daemon) or can be run via a cron job each minute

##### START: Test batch model
##IDs-BATCH:START#
#id1
#Id-2
#ID_3
#101
#other_ID
#Last-Id
##IDs-BATCH:END#
##### END: Test batch model

######################################## PERL MODULES

use strict;
use warnings;
use Cwd;
use Time::HiRes;
use Term::ANSIColor;

######################################## TERM COLORS

my $clr_error = ['bold yellow on_red'];
my $clr_fail = ['bold black on_bright_magenta'];
my $clr_notice = ['bold bright_white on_blue'];
my $clr_info = ['bold black'];
my $clr_hint = ['bold bright_yellow on_bright_blue'];
my $clr_msg_ok = ['green'];
my $clr_msg_yok = ['bold bright_green'];
my $clr_msg_xok = ['cyan'];

######################################## PARSE INI SETTINGS {{{SYNC-PERL-INIPARSE}}} ### DO NOT EDIT THIS SCRIPT !!! ### USE master.ini to store all settings ###

my %inisett = ();
my $cfname;
my $cfval;
open(CONF, "master.ini") || die "Failed to open config file master.ini";
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

# setting: max number of childs
my $max_threads = $inisett{'ChildThreads'}; # 3
my $batch_size = $inisett{'BatchSize'}; # 1000
my $offset_param = $inisett{'OffsetParam'}; # &ofs=

# setting: batch list URL
my $url_get = $inisett{'URLGetBatch'}; # http://some.url/?/action/get-batch

# setting: batch list URL Auth username [can be empty if no authentication is required]
my $user = $inisett{'URLAuthUser'}; # URL Auth User Name
# setting: batch list URL Auth password [can be empty if no authentication is required]
my $pass = $inisett{'URLAuthPassword'}; # URL Auth User Password

my $master_verbose = "yes";
if($inisett{'Verbose'} eq "no") {
	$master_verbose = "no";
}

######################################## RUNTIME

my $is_url = str_begins_with($url_get, 'http://') || str_begins_with($url_get, 'https://') || str_begins_with($url_get, 'ftp://');

my $num_args = $#ARGV + 1;
my $navOffset = -1;
my $fArg = "";

my $txtNavOffs = " [using BATCH-SIZE=".$batch_size."] ";

if($is_url eq "OK") {
	$url_get = $url_get.$batch_size;
}
if($num_args == 1) {
	$fArg = $ARGV[0];
	$fArg =~ s/^\s+|\s+$//g; #trim
	if($fArg =~ /^[0-9]*$/) { # test if in 0-9
		if($fArg >= 0) {
			$navOffset = $fArg;
			if($navOffset >= 0) {
				if($offset_param ne "") {
					$txtNavOffs .= " [using OFFSET=".($navOffset * $batch_size)." / PAGE-OFFSET=".$navOffset."] ";
					$txtNavOffs .= " [using OFFSET-PARAM=".$offset_param."] ";
					if($is_url eq "OK") {
						$url_get = $url_get.$offset_param.($navOffset * $batch_size);
					}
				}
			}
		}
	}
}

# Test if semaphore exists
my $dir = getcwd;
my $the_master_pid = $dir."/master-semaphore.pid";
if(-e $the_master_pid) {
	print colored($clr_error, "MASTER.ERR: Pid Semaphore already exist ... ".$the_master_pid);
	print "\n";
	exit;
}

# Catch CTRL+C
my $shutdown = "";
$SIG{INT} = sub { print colored($clr_notice, "MASTER.INF: STOP SIGNAL received ... please wait ..."); print "\n"; $shutdown = "stop"; };

# Create Semaphores Directory in the current directory of this script
system("mkdir -p ".$dir."/child-semaphores");
if(-d $dir."/child-semaphores") {
	# OK
} else {
	print colored($clr_error, "MASTER.ERR: Could Not Create Semaphores Dir: ".$dir."/child-semaphores");
	print "\n";
	exit;
}

# Cleanup Old Semaphores
my $expmis = 10;
print colored($clr_notice, "MASTER.INF: Cleaning up expired child pids, older than ${expmis} minutes ...");
print "\n";
my $cleanup = `find ./child-semaphores/ -type f -name "*.pid" -mmin +${expmis} -exec ls {} \\;`;
foreach my $cline (split /[\r\n]+/, $cleanup) {
	if($cline ne "") {
		if(-e $cline) {
			print "MASTER.INF: Removing expired child PID: ".$cline;
			print "\n";
			system("rm ".$cline); # remove the master semaphore if exists
		}
	}
}
print colored($clr_notice, "MASTER.INF: Initializing Data Pool ...");
print "\n";
sleep(5);

# Build the URL Request to get the Task List Batch
$url_get = str_single_quotes_escapeshellarg($url_get);
my $auth = "";
if($is_url eq "OK") {
	if($user ne "") {
		if($pass ne "") { # use auth: user and pass
			$auth = " --basic -u '".str_single_quotes_escapeshellarg($user).":".str_single_quotes_escapeshellarg($pass)."'";
		} else { # use auth: only user
			$auth = " --basic -u '".str_single_quotes_escapeshellarg($user)."'";
		}
	}
	if($auth eq "") {
		print colored($clr_notice,"############### Smart.Task.Engine // MASTER: Fetching URL".$txtNavOffs.": ".$url_get);
	} else {
		print colored($clr_notice, "############### Smart.Task.Engine // MASTER: Fetching URL [using AUTH=".$user."/*****"."]".$txtNavOffs.": ".$url_get);
	}
} else {
	if($url_get ne "") {
		print colored($clr_notice, "############### Smart.Task.Engine // MASTER: Fetching FILE: ".$url_get);
	} else {
		print colored($clr_error, "############### Smart.Task.Engine // MASTER: NO FILE OR URL TO FETCH !");
	}
}
print "\n";

# Get the Task List Batch from URL using curl
## G = get
## f = don't output if HTTP error
## s = silent, don't output warnings
## k = insecure, ignore self-signed certificate
my @arr_lines = ();
my $arr_len = 0;
my $response = "";

if($url_get ne "") {
	if($is_url eq "OK") {
		$response = `curl -Gfsk --connect-timeout 30 --max-time 600${auth} --url '${url_get}'`;
	} else {
		$response = `cat '${url_get}'`;
	}
} else {
	$response = "";
}
$response =~ s/^\s+|\s+$//g; #trim
if($response eq '') { # in perl == is only for numbers
	print colored($clr_fail, 'MASTER.ERR: EMPTY BATCH FILE ...');
	print "\n";
	exit;
}
#print $response."\n"; # Debug only: server response

# Test if server $response is empty
@arr_lines = split("\n", $response);
$arr_len = scalar @arr_lines;
if($arr_len <= 0) {
	print colored($clr_fail, 'MASTER.ERR: INVALID BATCH FILE LINES ...');
	print "\n";
	exit;
}

# Info: how many lines are in the server response
print colored($clr_info, "MASTER.INF: Server Answer batch have [".$arr_len."] lines");
print "\n";

# Trim $response
$arr_lines[0] =~ s/^\s+|\s+$//g;

# Test if server $response first line is: #IDs-BATCH:START#
if($arr_lines[0] ne "#IDs-BATCH:START#") {
	print colored($clr_fail, 'MASTER.ERR: INVALID BATCH FILE START-LINE ...');
	print "\n";
	exit;
}
#print "Checking 1st Line: OK"."\n";

# Test if server $response last line is: #IDs-BATCH:END#
if($arr_lines[$arr_len-1] ne "#IDs-BATCH:END#") {
	print colored($clr_fail, 'MASTER.ERR: INVALID BATCH FILE END-LINE ...');
	print "\n";
	exit;
}
#print "Checking last Line: OK"."\n";

# create the master semaphore
system("echo ".'MASTER'." > ".$the_master_pid);
if(-e $the_master_pid) {
	#OK
} else {
	print colored($clr_fail, 'MASTER.ERR: FAILED to Lock the Master Semaphore: ['.$the_master_pid.'] ...');
	print "\n";
	exit;
}

# Main Loop to launch Childs
if($arr_len > 2) {
	my $running_childs = 0;
	my $i = 0;
	while(($i < $arr_len) && ($shutdown eq "")) { # check also for the shutdown handler
		# trim each line
		$arr_lines[$i] =~ s/^\s+|\s+$//g;
		# skip first line
		if($arr_lines[$i] eq "#IDs-BATCH:START#") {
			$i++;
			next; # skip
		}
		# skip last line
		if($arr_lines[$i] eq "#IDs-BATCH:END#") {
			$i++;
			next; # skip
		}
		# test line to be a valid ID
		my $the_task_id = "";
		if($arr_lines[$i] =~ /^[_A-Za-z0-9\-]*$/) { # test if in A-Z a-z 0-9 _ -
			$the_task_id = $arr_lines[$i];
		}
		# count how many child semaphores are currently running
		$running_childs = sprintf("%d", `ls ./child-semaphores/ | wc -l`);
		# test semaphore ID line is non-empty (if empty and valid as tested with regex above)
		if($the_task_id eq "") { # if ID is empty
			print colored($clr_hint, "WARNING: INVALID Batch Line @ Loop #".$i." / PageOffset=".$navOffset." for [TaskID: ".$arr_lines[$i]."] :: Running Childs # [".$running_childs."] of max [".$max_threads."]");
			print "\n";
		} else { # if ID is non-empty
			print colored($clr_info, "@ Loop #".$i." / PageOffset=".$navOffset." for [TaskID: ".$the_task_id."] :: Running Childs # [".$running_childs."] of max [".$max_threads."]");
			print "\n";
			# build the file reference for the child semaphore to test if exists only (child pid semaphores are managed in child.pl)
			my $the_child_pid = "./child-semaphores/semaphore@".$the_task_id.".pid";
			if($running_childs < $max_threads) { # if number of semaphores is lower than allowed
				if(-e $the_child_pid) { # if this semaphore is running, SKIP (old behaviour was RETRY)
					#if($i < ($arr_len - 1)) { # fix: decrease only if not the last ID (doing so will jump next list)
					#	$i--;
					#}
					#print colored($clr_msg_xok, "@ Retry ... the SEMAPHORE is currently RUNNING [TaskID=".$the_task_id."]");
					if($master_verbose ne "no") {
						print colored($clr_msg_xok, "* SKIP Child ... the SEMAPHORE is currently RUNNING [TaskID=".$the_task_id."]");
						print "\n";
					}
					Time::HiRes::sleep(1);
				} else { # if semaphore is not running, launch new child
					if($master_verbose ne "no") {
						print colored($clr_msg_ok, "+ SPAWNING Child # ".($running_childs+1)." [TaskID=".$the_task_id."]");
						print "\n";
					}
					system("./child.pl ".$the_task_id." &");
					Time::HiRes::sleep(0.1);
				}
			} else { # if number of semaphores is higher than allowed
				$i--;
				if($master_verbose ne "no") {
					print colored($clr_msg_yok, "- WAITING Child ... to get a FREE SEMAPHORE ... Running Childs # ".$running_childs." [TaskID=".$the_task_id."]");
					print "\n";
				}
				Time::HiRes::sleep(1);
			}
		}
		$i++;
		Time::HiRes::sleep(0.1); # loop sleep 0.1 seconds to be able to catch CTRL+C
	}
} else {
	print colored($clr_hint, "INFO: NO TASKS AVAILABLE # [".$arr_len."]");
	print "\n";
}

print colored($clr_notice, "############### Smart.Task.Engine // MASTER: DONE");
print "\n";

# finally sleep 0.1 seconds to avoid overloading on high speed CPUs
Time::HiRes::sleep(1);

# remove the master.pid
if(-e $the_master_pid) {
	system("rm ".$the_master_pid); # remove the master semaphore if exists
}

######################################## INTERNAL FUNCTIONS {{{SYNC-PERL-FXS}}}

sub str_single_quotes_escapeshellarg {
	my $arg = shift;
	$arg =~ s/'/'\\''/g; # escape single quotes
	return "".$arg;
}

sub str_begins_with {
	my $ok = "";
	if(substr($_[0], 0, length($_[1])) eq $_[1]) {
		$ok = "OK";
	}
	return "".$ok;
}

######################################## EXIT

print colored($clr_notice, "[ MASTER.INF: SCRIPT TERMINATED ]");
print "\n";

exit 0;

#END
