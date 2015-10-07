//----------------------------------------------------------------------------
// registry.cpp
//
// registry manipulation stuff
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#include <iostream>
#include <algorithm>
#include <string.h>
#include "registry.h"
#include "error.h"
#include "util.h"

using std::string;

//----------------------------------------------------------------------------
// Create a RegPath object from a registry key, which must be either
// HKEY_CURRENT_USER or HKEY_LOCAL_MACHINE.
//----------------------------------------------------------------------------

RegPath :: RegPath( HKEY root ) : mRoot( root ), mPathKey( 0 ) {
	long res = 0;
	if ( mRoot == HKEY_CURRENT_USER ) {
		res = RegOpenKeyEx( mRoot, "Environment", 0, KEY_ALL_ACCESS, & mPathKey);
	}
	else if ( mRoot == HKEY_LOCAL_MACHINE ) {
		res = RegOpenKeyEx( mRoot,
				"SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Environment",
				0, KEY_ALL_ACCESS, & mPathKey);
	}
	else {
		throw Error( "Invalid root key in RegPath") ;
	}

	if ( res != ERROR_SUCCESS ) {
		throw Error( "Could not get registry key - " + LastWinError() );
	}

	const int BUFSIZE = 2048;
	BYTE buffer[ BUFSIZE + 1 ];
	DWORD bufflen = BUFSIZE;
	DWORD type = 0;
	res = RegQueryValueEx( mPathKey, "PATH", NULL, & type, buffer, & bufflen   );
	if ( res != ERROR_SUCCESS ) {
		throw Error( "Could not get registry value - " + LastWinError() );
	}
	buffer[ bufflen ] = 0;
	SplitPath( (const char *) buffer );
//	std::cout << "[" << buffer << "]" << std::endl;
}

//----------------------------------------------------------------------------
// Close the path key opened in ctor.
//----------------------------------------------------------------------------

RegPath :: ~RegPath() {
	if ( mPathKey ) {
		RegCloseKey( mPathKey );
	}
}

//----------------------------------------------------------------------------
// Get path as single semi-colon separated string
//----------------------------------------------------------------------------

string RegPath :: Str() const {
	string p;
	for( unsigned int i = 0; i < mPath.size(); i++ ) {
		if ( p != "" ) {
			p += ";";
		}
		p += mPath[i];
	}
	return p;
}

//----------------------------------------------------------------------------
// How many directories on path?
//----------------------------------------------------------------------------

unsigned int RegPath :: Count() const {
	return mPath.size();
}

//----------------------------------------------------------------------------
// Get zero-based directory
//----------------------------------------------------------------------------

string RegPath :: At( unsigned int  i) const {
	return mPath.at( i );
}

//----------------------------------------------------------------------------
// Helper to split path at ';' character
//----------------------------------------------------------------------------

void RegPath :: SplitPath( const std::string & path ) {
	string::size_type pos = 0;

	while( pos != string::npos ) {
		string::size_type fpos = path.find( ";", pos );
		string s;
		if ( fpos == string::npos ) {
			s = path.substr( pos );
			pos = string::npos;
		}
		else {
			s = path.substr( pos, fpos - pos );
			pos = fpos + 1;
		}
		if ( s.find_first_not_of( " \t" ) != string::npos ) {
			mPath.push_back( s );
		}
	}
}

//----------------------------------------------------------------------------
// As stricmp seems to have disappeared....
//----------------------------------------------------------------------------

static bool Same( const std::string & a, const std::string &  b ) {
	if ( a.size() == b.size() ) {
		for ( unsigned int i = 0; i < a.size(); i++ ) {
			if ( toupper( a[i] ) != toupper( b[i] )) {
				return false;
			}
		}
		return true;
	}
	else {
		return false;
	}
}

//----------------------------------------------------------------------------
// See if path contains adir
//----------------------------------------------------------------------------

bool RegPath :: Find( const string & adir, Env ev ) const {

	for ( unsigned int i = 0; i < mPath.size(); i++ ) {
		string e = ev == Expand ? ExpandPath( mPath[i] ) : mPath[i];
		if ( Same( e, adir ) ) {
			return true;
		}
	}
	return false;
}

//----------------------------------------------------------------------------
// Add directory to path - no check for multiple adds of same directory
//----------------------------------------------------------------------------

bool RegPath :: Add( const string & adir ) {
	mPath.push_back( adir );
	UpdateReg();
	return true;
}

//----------------------------------------------------------------------------
// Helper to update registry with current path.
//----------------------------------------------------------------------------

void RegPath :: UpdateReg() {
	string newpath;
	for ( unsigned int i = 0; i < mPath.size(); i++ ) {
		newpath += mPath[i] + ";";
	}

	long res = RegSetValueEx( mPathKey, "PATH", 0, REG_EXPAND_SZ,
								(BYTE *)newpath.c_str(), newpath.size() + 1 );
	if ( res != ERROR_SUCCESS ) {
		throw Error( "Could not add update path in registry - " + LastWinError() );
	}
	//NotifyChanges();
}


//----------------------------------------------------------------------------
// Find string in path, ignoring case
//----------------------------------------------------------------------------

bool RegPath :: RemoveIC( const std::string & adir ) {
	for( VecType::iterator it = mPath.begin(); it != mPath.end(); ++ it ) {
		if ( Same( adir, * it ) ) {
			mPath.erase( it );
			return true;
		}
	}
	return false;
}

//----------------------------------------------------------------------------
// Remove single instance of adir from path, updating registry.
//----------------------------------------------------------------------------

bool RegPath :: Remove( const string & adir ) {
	if ( RemoveIC( adir ) ) {
		UpdateReg();
		return true;
	}
	else {
		return false;
	}
}

//----------------------------------------------------------------------------
// Replace entire path
//----------------------------------------------------------------------------

void RegPath :: ReplaceAll( const string & apath ) {
	long res = RegSetValueEx( mPathKey, "PATH", 0, REG_EXPAND_SZ,
								(BYTE *) apath.c_str(), apath.size() + 1 );
	if ( res != ERROR_SUCCESS ) {
		throw Error( "Could not add update path in registry - " + LastWinError() );
	}
}

// end


