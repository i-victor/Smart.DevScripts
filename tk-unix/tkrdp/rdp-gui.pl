#!/usr/bin/perl
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.

# v.159908 (c) unix-world.org

use warnings;
use File::Basename;

use Tk;
use Tk::LabEntry;


my $window = MainWindow->new();
$window->minsize(qw(550 250));
$window->maxsize(qw(550 250));
$window->title("  RDP Client  ");

$window->Label(-text => ' ')->pack();

my $frm_01 = $window->Frame(-borderwidth => 1, -relief => 'ridge')->pack(-ipadx => 5, -ipady => 5);
my $fld_srv = $frm_01->LabEntry(-label => 'Server:', -textvariable => \$servername, -width => 20, -labelPack => [-side => 'left', -anchor => 'n'], -relief => 'sunken')->pack(-side => 'left');
$frm_01->Label(-text, '     ')->pack(-side => 'left');
my $fld_dom = $frm_01->LabEntry(-label=>'Domain:', -textvariable => \$serverdomain, -width => 20, -labelPack => [-side => 'left', -anchor => 'n'], -relief => 'sunken')->pack(-side => 'left');

$window->Label(-text => ' ')->pack();

my $frm_02 = $window->Frame(-borderwidth => 1, -relief => 'ridge')->pack(-ipadx => 5, -ipady => 5);
my $fld_usr = $frm_02->LabEntry(-label => 'Username:', -textvariable => \$username, -width => 20, -labelPack => [-side => 'left', -anchor => 'n'], -relief => 'sunken')->pack(-side => 'left');
$frm_02->Label(-text, '     ')->pack(-side => 'left');
my $fld_pass = $frm_02->LabEntry(-label => 'Password:', -show => '*', -textvariable => \$userpass, -width => 15, -labelPack => [-side => 'left', -anchor => 'n'], -relief => 'sunken')->pack(-side => 'left');

$window->Label(-text => ' ')->pack();

my $frm_04 = $window->Frame(-borderwidth => 1, -relief => 'ridge')->pack(-ipadx => 5, -ipady => 5);
$frm_04->Label(-text, 'Info: RDP will start in FullScreen Mode. Use CTRL+Alt+ENTER to toggle out.')->pack(-side => 'bottom');

my $frm_03 = $window->Frame(-borderwidth, 1, -relief, 'flat')->pack(-side, 'bottom' , -ipadx, 10, -ipady, 10);
my $btn_conx = $frm_03->Button(-text => "Connect", -width => 18, -command => \&do_connection)->pack(-side => 'left');
$frm_03->Label(-text, '          ')->pack(-side => 'left');
my $btn_exit = $frm_03->Button(-text => "Quit", -width => 18, -command => \&do_exit)->pack(-side => 'left');

#####

MainLoop();

#####

sub do_connection {
	#
	if(defined($servername)) {
		#
		$fld_srv->configure(-state => 'disabled');
		$fld_dom->configure(-state => 'disabled');
		$fld_usr->configure(-state => 'disabled');
		$fld_pass->configure(-state => 'disabled');
		$btn_conx->configure(-text => "Connecting ...", -state => 'disabled');
		# the rdesktop requires at least server or server:port
		$cmd = "/usr/local/bin/rdesktop -f";
		# domain
		if(defined($serverdomain)) {
			$cmd .= " -d ".&escapeshellarg($serverdomain);
		}
		# username
		if(defined($username)) {
			$cmd .= " -u ".&escapeshellarg($username);
		}
		# password
		if(defined($userpass)) {
			$cmd .= " -p ".&escapeshellarg($userpass);
		}
		# run the command
		$cmd .= " ".&escapeshellarg($servername);
		#print "opening connection ...\n".$cmd."\n";
		system $cmd;
		#
	}
}

sub do_exit {
	print "... closing connection.\n";
	print "";
	exit;
}

#####

sub escapeshellarg {
    my $arg = shift;
    $arg =~ s/'/'\\''/g;
    return "'".$arg."'";
}

#####

#END
