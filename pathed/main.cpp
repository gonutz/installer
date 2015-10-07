//----------------------------------------------------------------------------
// main.cpp
//
// main for pathed
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#include <iostream>
#include <set>
#include <direct.h>  // for getcwd()
#include "cmdline.h"
#include "error.h"
#include "registry.h"
#include "util.h"

using namespace std;

//----------------------------------------------------------------------------
// Command and option names
//----------------------------------------------------------------------------

enum FlagName { fnNone, fnAdd, fnRemove, fnForce, fnGrep,
				fnQuery, fnVerify, fnPrune, fnList, fnSys,
				fnExpand, fnEnv, fnUnix, fnCwd, fnUnCwd };

//----------------------------------------------------------------------------
// Globals set during command parsing
//----------------------------------------------------------------------------

FlagName CommandName = fnNone;
bool Expand = false, CheckExist = true, UseSys = false, Unix = false;
string CommandParam = "";

//----------------------------------------------------------------------------
// Lookup commands and options via string
//----------------------------------------------------------------------------

struct Flag {
	FlagName mName;
	const char * const mShort;
	const char * const mLong;
	bool mCmd;						// is this a command?
	int mParamCount;				// if so, how many params?
};

Flag CmdLineFlags[] = {
	{ fnAdd, 		"-a", "--add", true, 1 },
	{ fnCwd, 		"-c", "--cwd", true, 0 },
	{ fnUnCwd, 		"-d", "--uncwd", true, 0 },
	{ fnRemove, 	"-r", "--remove", true, 1 },
	{ fnList, 		"-l", "--list", true, 0 },
	{ fnQuery, 		"-q", "--query", true, 1 },
	{ fnVerify, 	"-v", "--verify", true, 0 },
	{ fnPrune, 		"-p", "--prune", true, 0 },
	{ fnSys, 		"-s", "--system", false, 0 },
	{ fnExpand, 	"-x", "--expand", false, 0 },
	{ fnForce, 		"-f", "--force", false, 0 },
	{ fnGrep, 		"-g", "--grep", true, 1 },
	{ fnEnv,		"-e", "--env", true, 0 },
	{ fnUnix,		"-u", "--unix", false, 0 },

	{ fnNone, NULL, NULL, false, 0 }		// must be last
};


//----------------------------------------------------------------------------
// Get list of commands, separated by commas
//----------------------------------------------------------------------------

string CommandList() {
	int i = 0;
	string cl;
	while( CmdLineFlags[i].mName != fnNone ) {
		if ( CmdLineFlags[i].mCmd ) {
			if ( cl != "" ) {
				cl += ", ";
			}
			cl +=CmdLineFlags[i].mShort;
		}
		i++;
	}
	return cl;
}

//----------------------------------------------------------------------------
// Get flag name from string rep.
//----------------------------------------------------------------------------

FlagName StringToFlag( const string & s  ) {
	int i = 0;
	while( CmdLineFlags[i].mName != fnNone ) {
		if ( s == CmdLineFlags[i].mShort || s == CmdLineFlags[i].mLong ) {
			return CmdLineFlags[i].mName;
		}
		i++;
	}
	return fnNone;
}

//----------------------------------------------------------------------------
// Get values associated with an option
//----------------------------------------------------------------------------

pair <bool,int> GetFlagValues( FlagName f ) {
	int i = 0;
	while( CmdLineFlags[i].mName != fnNone ) {
		if (CmdLineFlags[i].mName == f ) {
			return make_pair( CmdLineFlags[i].mCmd, CmdLineFlags[i].mParamCount );
		}
		i++;
	}
	throw "wonky flag ";	// never happen
}

//----------------------------------------------------------------------------
// Parse supplied command line, setting globals.
//----------------------------------------------------------------------------

void ParseCommandLine( CmdLine & cl  ) {
	while( cl.Argc() > 1 ) {
		string sflag = cl.Argv(1);
		cl.Shift();
		FlagName fn = StringToFlag( sflag );
		if ( fn == fnNone ) {
			throw Error( "Invalid command line option: " + sflag );
		}
		pair <bool,int> vals = GetFlagValues( fn );
		if ( vals.second ) {
			if ( cl.Argc() == 1 ) {
				throw Error( "Missing command line parameter for " + sflag );
			}
			CommandParam = cl.Argv(1);
			cl.Shift();
		}
		if ( vals.first ) {
			if ( CommandName != fnNone ) {
				throw Error( "Only one command option allowed" );
			}
			CommandName = fn;
		}
		else {
			switch( fn ) {
				case fnExpand:		Expand = true; break;
				case fnForce:		CheckExist = false ; break;
				case fnSys:			UseSys  = true; break;
				case fnUnix:		Unix = true; break;
				default:			throw Error( "bad option " );
			}
		}
	}
	if ( CommandName == fnNone ) {
		throw Error( "Need one of " + CommandList() );
	}
}

//----------------------------------------------------------------------------
// If the global Unix flag is set, replace all backslashes in path with
// forward slashes, returning new path.
//----------------------------------------------------------------------------

string ConvertSep( const string & path ) {
	if ( ! Unix ) {
		return path;
	}
	string np;
	for( unsigned int i = 0; i < path.size(); i++ ) {
		if ( path[i] == '\\' ) {
			np += '/';
		}
		else {
			np += path[i];
		}
	}
	return np;
}

//----------------------------------------------------------------------------
// List PATH to stdout, one directory per line
//----------------------------------------------------------------------------

void ListPath() {
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	for ( unsigned int i = 0; i < path.Count(); i++ ) {
		cout << ( Expand
					? ConvertSep( ExpandPath( path.At(i) ) )
					: ConvertSep( path.At(i) ) ) << "\n";
	}
}

//----------------------------------------------------------------------------
// Add an entry to the path
//----------------------------------------------------------------------------

void AddPath( const char * cwd = NULL ) {
	string pathstr = cwd == NULL ? CommandParam : cwd;
	if ( pathstr == ""  ) {
		throw Error( "Need directory to add" );
	}
	if ( CheckExist ) {
		DWORD attr = GetFileAttributes( pathstr.c_str() );
		if ( attr == INVALID_FILE_ATTRIBUTES || ! (attr & FILE_ATTRIBUTE_DIRECTORY ) ) {
			throw Error( "No such directory: " + pathstr );
		}
	}

	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	if ( path.Find( pathstr, RegPath::NoExpand ) ) {
		return;
	}
	path.Add( pathstr );
}

//----------------------------------------------------------------------------
// add current directory to path by calling addpath
//----------------------------------------------------------------------------

void AddCwd() {
	const int BUFSIZE = 2048;
	char buffer[BUFSIZE];
	char * p = getcwd( buffer, BUFSIZE );
	if ( p == NULL ) {
		throw Error( "Could not get working directory name" );
	}
	AddPath( buffer );
}

//----------------------------------------------------------------------------
// Remove entry from the path
//----------------------------------------------------------------------------

void RemovePath( const char * cwd = NULL ) {
	string pathstr = cwd == NULL ? CommandParam : cwd;
	if ( pathstr == "" ) {
		throw Error( "Need directory to remove" );
	}
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	if ( ! path.Find( pathstr , RegPath::NoExpand ) ) {
		throw Error( pathstr + " is not on the path" );
	}
	if ( ! path.Remove( pathstr ) ) {
		throw Error( pathstr + "not found on path" );
	}
}

//----------------------------------------------------------------------------
// Remove  current directory from  path by
//----------------------------------------------------------------------------

void UnCwd() {
	const int BUFSIZE = 2048;
	char buffer[BUFSIZE];
	char * p = getcwd( buffer, BUFSIZE );
	if ( p == NULL ) {
		throw Error( "Could not get working directory name" );
	}
	RemovePath( buffer );
}

//----------------------------------------------------------------------------
// Prune duplicates and non-existent dirs from path. Use a set to detect
// dupes, but actually work with vector to meaintain path order.
//----------------------------------------------------------------------------

void PrunePath() {
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );

	typedef std::set <string> DirSet;
	DirSet uniq;
	std::vector <string> ordered;

	for ( unsigned int i = 0; i < path.Count(); i++ ) {
		string dir = path.At( i );
		std::pair<DirSet::iterator, bool> ok = uniq.insert( dir );
		if ( ok.second ) {
			ordered.push_back( dir );
		}
		else {
			cout << "Pruned: " << dir << endl;
		}
	}

	string entry;

	for ( unsigned int i = 0; i < ordered.size(); i++ ) {
		string dir = ExpandPath( ordered[i] );
		DWORD attr = GetFileAttributes( dir.c_str() );
		if ( attr == INVALID_FILE_ATTRIBUTES ||
						! (attr & FILE_ATTRIBUTE_DIRECTORY ) ) {
			cout << "Pruned: " << ConvertSep( ordered[i] ) << endl;
		}
		else {
			if ( entry != "" ) {
				entry += ";";
			}
			entry += ordered[i];
		}

	}
	path.ReplaceAll( entry );
	NotifyChanges();
	std::cout << ConvertSep( entry ) << std::endl;
}

//----------------------------------------------------------------------------
// See if directory is on the path, if so return success code (not boolean!)
//----------------------------------------------------------------------------

int FindPath() {
	if ( CommandParam == "" ) {
		throw Error( "Need directory name" );
	}
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	return  path.Find( CommandParam, Expand ? RegPath::Expand : RegPath::NoExpand ) ? 0 : 1;
}

//----------------------------------------------------------------------------
// Search path for file
//----------------------------------------------------------------------------

int GrepPath() {
	if ( CommandParam == "" ) {
		throw Error( "Need file name" );
	}
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	int found = 0;
	for ( unsigned int i = 0; i < path.Count(); i++ ) {
		string epath = ExpandPath( path.At(i) );
		if ( epath == "" || epath[epath.size()-1] != '\\' ) {
			epath += '\\';
		}
		epath += CommandParam;
		// out << epath << endl;
		DWORD attr = GetFileAttributes( epath.c_str() );
		if ( attr != INVALID_FILE_ATTRIBUTES ) {
			found++;
			cout << ConvertSep( epath ) << endl;
		}
	}
	return found == 0 ? 1 : 0;
}

//----------------------------------------------------------------------------
// List path in PATH environment variable - flags except -u have no effect.
//----------------------------------------------------------------------------

int EnvPath() {
	const char * p = getenv( "PATH" );
	if ( p == 0 ) {
		throw Error( "No PATH variable in environment!" );
	}
	while( * p ) {
		if ( * p == ';' ) {
			cout << '\n';
		}
		else if ( * p == '\\' && Unix ) {
			cout << "/";
		}
		else {
			cout << * p;
		}
		p++;
	}
	cout << endl;
	return 0;
}

//----------------------------------------------------------------------------
// Verify directories on path exist.
//----------------------------------------------------------------------------

int VerifyPath() {
	RegPath path( UseSys ? HKEY_LOCAL_MACHINE : HKEY_CURRENT_USER );
	int bad = 0;
	for ( unsigned int i = 0; i < path.Count(); i++ ) {
		string epath = ExpandPath( path.At(i) );
		DWORD attr = GetFileAttributes( epath.c_str() );
		if ( attr == INVALID_FILE_ATTRIBUTES ) {
			cout << "No such directory: " << epath << "\n";
			bad++;
		}
		else if ( ! (attr & FILE_ATTRIBUTE_DIRECTORY ) ) {
			cout << "Not a directory: " << epath << "\n";
			bad++;
		}
	}
	return bad == 0 ? 0 : 1;
}


//----------------------------------------------------------------------------
// Display help
//----------------------------------------------------------------------------

void Help() {

	cout <<

	"\npathed is a command-line tool for changing and querying the path in the registry\n\n"
	"Version 0.9\n"
	"Copyright (C) 2012 Neil Butterworth\n\n"
	"usage: pathed [-a dir | -r dir | -c | -d | -e | -l | -q dir | -v | -p | -g file] [-s] [-f] [-x] [-u] \n\n"
	"pathed -a dir    adds dir to the path in  the registry\n"
	"pathed -c        add current working directory to path in registry\n"
	"pathed -d        remove current working directory from path in registry\n"
	"pathed -r dir    removes  dir from the path in the registry\n"
	"pathed -l        lists the entries on the current path in the registry\n"
	"pathed -e        lists the entries on the current path in the PATH environment variable\n"
	"pathed -q dir    queries registry, returns 0 if dir is on path, 1 otherwise\n"
	"pathed -g file   searches (greps) the path for all occurrences of file\n"
	"pathed -v        verifies that all directories on the path exist\n"
	"pathed -p        prunes the path by removing duplicates and non-existent directories\n\n"
	"By default, pathed works on the path in HKEY_CURRENT_USER. You can make it use\n"
	"the system path in HKEY_LOCAL_MACHINE by using the -s flag.\n\n"
	"Normally pathed will check a directory exists on disk before adding it to the\n"
	"path. To prevent this, use the -f flag.\n\n"
	"Paths containing environment variables such as %systemroot% will not normally have\n"
	"the variables expanded to their values. To expand them, use the -x flag\n\n"
	"On output only, use the -u flag to produce UNIX-style paths, using forward-slash\n"
	"as the separator\n\n"
	"AS WITH ALL COMMANDS THAT CHANGE THE REGISTRY, PATHED CAN CAUSE DAMAGE IF YOU\n"
	"DO NOT KNOW WHAT YOU ARE DOING. IF IN DOUBT, DO NOT USE IT!\n"

	<< endl;
}

//----------------------------------------------------------------------------
// Main for pathed
//----------------------------------------------------------------------------

int main( int argc, char *argv[] ) {
	try {
		CmdLine cl( argc, argv );
		if ( cl.Argc() == 1 ) {
			Help();
			return 0;
		}

		ParseCommandLine( cl );

		switch( CommandName ) {
			case fnAdd:		AddPath(); break;
			case fnRemove:	RemovePath(); break;
			case fnQuery:	return FindPath();
			case fnList:	ListPath(); break;
			case fnVerify:	return VerifyPath(); break;
			case fnPrune:	PrunePath(); break;
			case fnGrep:	GrepPath(); break;
			case fnEnv:	    EnvPath(); break;
			case fnCwd:		AddCwd(); break;
			case fnUnCwd:	UnCwd(); break;
			default:		throw Error( "bad command switch" );
		}

		return 0;
	}
	catch( const Error & e ) {
		cerr << e.what() << endl;
		return 1;
	}
	catch( ... ) {
		cerr << "Unexpected exception" << endl;
		return 1;
	}
}
