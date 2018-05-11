#5/11/2018
I decided to update this code to fix a problem where an existing query would error out and not be added to the pack. This revision catches the 409 http.statusCode
and performs a api call to get all the queries in the db, it then iterates through the json query objects until it finds a match for the name of the existing query
it then returns that existing query's Query ID so that the script can continue to add the existing query to the pack. For our testing we used the Palantir OSquery 
packs as the bulk import files. 

WARNING: THIS script must be run from the directory which contains the packs but the code itself cannot be in the packs directory. 

WARNING: Also I never used GOLANG before this so the efficacy of this patch is not under any kind of warranty

Original Code: https://gist.github.com/marpaia/9e061f81fa60b2825f4b6bb8e0cd2c77

Palantir Packs used for testing: https://github.com/palantir/osquery-configuration/tree/master/Endpoints/packs 

#################################################################################################################################################################

# Query Pack Import Tool

To run the tool, download the `import.go` file somewhere locally. It can then be executed via `go run`:

```
$ go run ./import.go -help

Usage of /var/folders/wp/6fkmvjf11gv18tdprv4g2mk40000gn/T/go-build234469651/command-line-arguments/_obj/exe/import:
  -hostname string
    	Kolide server hostname (default "https://localhost:8080")
  -pack_dir string
    	Directory of packs
  -token string
    	Kolide authentication token
exit status 2
```

## Usage

The `import.go` script accepts three command-line flags:

#### hostname

This is the hostname of your Kolide server. This should be in the format `https://foobar.xyz:1234`. Note that `https://` must be prepended to the hostname and there should be no trailing slashes.

#### pack_dir

This is the local directory where all of your packs are located. Absolute paths are preferred.

#### token

This is a valid authentication token, so that the tool can communicate with the Kolide API. This isn't super elegant, but to get this token:

- go to Kolide
- open the web inspector
- select the "XHR" tab
- refresh the page if you don't see any requests in the sidebar
- select a request from the sidebar
- under the headers tab, look for "Request Headers"
- copy JUST the token (NOT the "Bearer " text)

![use the force](http://i.imgur.com/6iZhOs5.png)

## Sample output

To see what this script should output on a successful execution, observe the following output:

```
$ go run import.go -pack_dir ~/Desktop -token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzZXNzaW9uX2tleSI6InQxWEdqU1JvU0N3NWNubExIcWduSSszeDNOWmFDMjNXalRzMm5kVU44VnVRb3NUem5MNzI5OWhObFBsREFuK09JenBCcWdNaWY1TGtRRUV3Q3lVdHNBPT0ifQ.XBxgopLiuwUSQX_NnHmnX4A0oVDVFk9tgmwuyAC7IGQ"
2017/03/07 11:28:47 .DS_Store is not a query pack... skipping.
2017/03/07 11:28:47 .localized is not a query pack... skipping.
2017/03/07 11:28:47 Created pack centos.conf (30)
2017/03/07 11:28:47 Created query listening_ports (276)
2017/03/07 11:28:47 Added query listening_ports to pack centos.conf (35)
2017/03/07 11:28:47 Created query disk_encryption (277)
2017/03/07 11:28:47 Added query disk_encryption to pack centos.conf (36)
2017/03/07 11:28:47 Created query schedule (278)
2017/03/07 11:28:47 Added query schedule to pack centos.conf (37)
2017/03/07 11:28:47 Created query events (279)
2017/03/07 11:28:47 Added query events to pack centos.conf (38)
2017/03/07 11:28:47 Created query sudoers (280)
2017/03/07 11:28:47 Added query sudoers to pack centos.conf (39)
2017/03/07 11:28:47 Created query etc_hosts (281)
2017/03/07 11:28:47 Added query etc_hosts to pack centos.conf (40)
2017/03/07 11:28:47 Created query kernel_modules (282)
2017/03/07 11:28:47 Added query kernel_modules to pack centos.conf (41)
2017/03/07 11:28:47 Created query ramdisk (283)
2017/03/07 11:28:47 Added query ramdisk to pack centos.conf (42)
2017/03/07 11:28:47 Created query authorized_keys (284)
2017/03/07 11:28:47 Added query authorized_keys to pack centos.conf (43)
2017/03/07 11:28:47 Created query mounts (285)
2017/03/07 11:28:47 Added query mounts to pack centos.conf (44)
2017/03/07 11:28:47 Created query os_version (286)
2017/03/07 11:28:47 Added query os_version to pack centos.conf (45)
2017/03/07 11:28:47 Created query logged_in_users (287)
2017/03/07 11:28:47 Added query logged_in_users to pack centos.conf (46)
2017/03/07 11:28:47 Created query iptables (288)
2017/03/07 11:28:47 Added query iptables to pack centos.conf (47)
2017/03/07 11:28:47 Created query root_directory (289)
2017/03/07 11:28:47 Added query root_directory to pack centos.conf (48)
2017/03/07 11:28:47 Created query private_ssh_key (290)
2017/03/07 11:28:47 Added query private_ssh_key to pack centos.conf (49)
2017/03/07 11:28:47 Created query acpi_tables (291)
2017/03/07 11:28:47 Added query acpi_tables to pack centos.conf (50)
2017/03/07 11:28:47 Created query crontab (292)
2017/03/07 11:28:47 Added query crontab to pack centos.conf (51)
2017/03/07 11:28:47 Created query open_sockets (293)
2017/03/07 11:28:47 Added query open_sockets to pack centos.conf (52)
2017/03/07 11:28:47 Created query ip_forwarding (294)
2017/03/07 11:28:47 Added query ip_forwarding to pack centos.conf (53)
2017/03/07 11:28:47 Created query suid_bin (295)
2017/03/07 11:28:47 Added query suid_bin to pack centos.conf (54)
2017/03/07 11:28:47 Created query arp_cache (296)
2017/03/07 11:28:47 Added query arp_cache to pack centos.conf (55)
2017/03/07 11:28:47 Created query osquery_info (297)
2017/03/07 11:28:47 Added query osquery_info to pack centos.conf (56)
2017/03/07 11:28:47 Created query rpm_packages (298)
2017/03/07 11:28:47 Added query rpm_packages to pack centos.conf (57)
2017/03/07 11:28:47 Created query kernel_info (299)
2017/03/07 11:28:47 Added query kernel_info to pack centos.conf (58)
2017/03/07 11:28:47 Created query usb_devices (300)
2017/03/07 11:28:47 Added query usb_devices to pack centos.conf (59)
2017/03/07 11:28:47 Created query last (301)
2017/03/07 11:28:47 Added query last to pack centos.conf (60)
2017/03/07 11:28:47 import.go is not a query pack... skipping.
2017/03/07 11:28:47 import_backup.go is not a query pack... skipping.
```

This was developed with the following as the objective packs file:

```json
{
  "queries": {
    "acpi_tables": {
      "query": "select * from acpi_tables;",
      "interval": 86400,
      "description": "General reporting and heuristics monitoring."
    },
    "kernel_info": {
      "query": "select * from kernel_info;",
      "interval": 7200,
      "description": "Report the booted kernel, potential arguments, and the device."
    },
    "usb_devices": {
      "query": "select * from usb_devices;",
      "interval": 7200,
      "description": "Report an inventory of USB devices. Attaches and detaches will show up in hardware_events."
    },
    "crontab": {
      "query" : "select * from crontab join hash using (path);",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves all the jobs scheduled in crontab in the target system.",
      "value" : "Identify malware that uses this persistence mechanism to launch at a given interval"
    },
    "etc_hosts": {
      "query" : "select * from etc_hosts;",
      "interval" : "86400",
      "version" : "1.4.5",
      "description" : "Retrieves all the entries in the target system /etc/hosts file.",
      "value" : "Identify network communications that are being redirected. Example: identify if security logging has been disabled"
    },
    "kernel_modules": {
      "query" : "select * from kernel_modules;",
      "interval" : "3600",
      "platform" : "linux",
      "version" : "1.4.5",
      "description" : "Retrieves all the information for the current kernel modules in the target Linux system.",
      "value" : "Identify malware that has a kernel module component."
    },
    "last": {
      "query" : "select * from last;",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves the list of the latest logins with PID, username and timestamp.",
      "value" : "Useful for intrusion detection and incident response. Verify assumptions of what accounts should be accessing what systems and identify machines accessed during a compromise."
    },
    "open_sockets": {
      "query" : "select distinct pid, family, protocol, local_address, local_port, remote_address, remote_port, path from process_open_sockets where path <> '' or remote_address <> '';",
      "interval" : "86400",
      "version" : "1.4.5",
      "description" : "Retrieves all the open sockets per process in the target system.",
      "value" : "Identify malware via connections to known bad IP addresses as well as odd local or remote port bindings"
    },
    "logged_in_users": {
      "query" : "select liu.*, p.name, p.cmdline, p.cwd, p.root from logged_in_users liu, processes p where liu.pid = p.pid;",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves the list of all the currently logged in users in the target system.",
      "value" : "Useful for intrusion detection and incident response. Verify assumptions of what accounts should be accessing what systems and identify machines accessed during a compromise."
    },
    "ip_forwarding": {
      "query" : "select * from system_controls where oid = '4.30.41.1' or oid = '4.2.0.1';",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves the current status of IP/IPv6 forwarding.",
      "value" : "Identify if a machine is being used as relay."
    },
    "mounts": {
      "query" : "select * from mounts;",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves the current list of mounted drives in the target system.",
      "value" : "Scope for lateral movement. Potential exfiltration locations. Potential dormant backdoors."
    },
    "ramdisk": {
      "query" : "select * from block_devices where type = 'Virtual Interface';",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves all the ramdisk currently mounted in the target system.",
      "value" : "Identify if an attacker is using temporary, memory storage to avoid touching disk for anti-forensics purposes"
    },
    "listening_ports": {
      "query" : "select name, path, listening_ports.* from processes join listening_ports using (pid);",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves all the listening ports in the target system.",
      "value" : "Detect if a listening port iis not mapped to a known process. Find backdoors."
    },
    "suid_bin": {
      "query" : "select suid_bin.*, md5, sha1, sha256 from hash join suid_bin using (path);",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves all the files in the target system that are setuid enabled.",
      "value" : "Detect backdoor binaries (attacker may drop a copy of /bin/sh). Find potential elevation points / vulnerabilities in the standard build."
    },
    "arp_cache": {
      "query" : "select * from arp_cache;",
      "interval" : "3600",
      "version" : "1.4.5",
      "description" : "Retrieves the ARP cache values in the target system.",
      "value" : "Determine if MITM in progress."
    },
    "disk_encryption": {
      "query" : "select * from disk_encryption;",
      "interval" : "86400",
      "version" : "1.4.5",
      "description" : "Retrieves the current disk encryption status for the target system.",
      "value" : "Identifies a system potentially vulnerable to disk cloning."
    },
    "iptables": {
      "query" : "select * from iptables;",
      "interval" : "3600",
      "platform" : "linux",
      "version" : "1.4.5",
      "description" : "Retrieves the current filters and chains per filter in the target system.",
      "value" : "Verify firewall settings are as restrictive as you need. Identify unwanted firewall holes made by malware or humans"
    },
    "osquery_info": {
      "query" : "select * from time, osquery_info;",
      "interval" : "86400",
      "version" : "1.4.5",
      "description" : "Retrieves the current version of the running osquery in the target system and where the configuration was loaded from.",
      "value" : "Identify if your infrastructure is running the correct osquery version and which hosts may have drifted"
    },
    "os_version": {
      "query" : "select * from os_version;",
      "interval" : "86400",
      "version" : "1.4.5",
      "description" : "Retrieves information from the Operative System where osquery is currently running.",
      "value" : "Identify out of date operating systems or version drift across your infrastructure"
    },
    "rpm_packages": {
      "query" : "select * from rpm_packages;",
      "interval" : "86400",
      "platform" : "redhat,centos",
      "version" : "1.4.5",
      "description" : "Retrieves all the installed RPM packages in the target Linux system.",
      "value" : "General security posture."
    },
    "schedule": {
      "query": "select name, interval, executions, output_size, wall_time, (user_time/executions) as avg_user_time, (system_time/executions) as avg_system_time, average_memory, last_executed from osquery_schedule;",
      "interval": 7200,
      "removed": false,
      "version": "1.6.0",
      "description": "Report performance for every query within packs and the general schedule."
    },
    "events": {
      "query": "select name, publisher, type, subscriptions, events, active from osquery_events;",
      "interval": 86400,
      "removed": false,
      "description": "Report event publisher health and track event counters."
    },
    "sudoers": {
      "query": "select file.*, md5, sha1, sha256 from hash join file using (path) where path = '/etc/sudoers';",
      "interval": 86400,
      "removed": false,
      "description": "Rules for running commands as other users via sudo"
    },
    "authorized_keys": {
      "query": "select authorized_keys.* from users join authorized_keys using (uid);",
      "interval": 86400,
      "removed": false,
      "description": "A line-delimited authorized_keys table"
    },
    "root_directory": {
      "query": "select file.* from users join file using (uid) where path = '/' or path = '/root';",
      "interval": 86400,
      "removed": false,
      "description": "Information on the permissions of the root directory"
    },
    "private_ssh_key": {
      "query": "select file.* from users join file using (uid) where path like '/home/%%/.ssh/id_rsa' or path like '/home/%%/.ssh/id_dsa' or path like '/root/.ssh/id_rsa' or path like '/root/.ssh/id_dsa' or path like '/data/home/%%/.ssh/id_rsa' or path like '/data/home/%%/.ssh/id_dsa';",
      "interval": 86400,
      "removed": false,
      "description": "Ensure that users do not have private ssh keys on servers"
    }
  }
}
```
