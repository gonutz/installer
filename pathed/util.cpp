//----------------------------------------------------------------------------
// util.cpp
//
// utility functions
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#include "util.h"
#include <cstdlib>
#include <windows.h>
using std::string;

//----------------------------------------------------------------------------
// Get env var value
//----------------------------------------------------------------------------

string GetEnv( const string &  name ) {
	const char * val = getenv( name.c_str() );
	return val == 0 ? "" : val;
}

//----------------------------------------------------------------------------
// Expand occurrences of %X% in the path name with the relevant environment
// variable setting.
//----------------------------------------------------------------------------

string ExpandPath( const std::string &  adir ) {
	string rv, envname;
	bool inenv = false;
	unsigned int i = 0;
	while( i < adir.size() ) {
		char c = adir[i++];
		if ( c == '%' && inenv ) {
			rv += GetEnv( envname );
			inenv = false;
			envname = "";
		}
		else if ( c == '%' && ! inenv ) {
			envname = "";
			inenv = true;
		}
		else if ( inenv ) {
			envname += c;
		}
		else {
			rv += c;
		}
	}
	if ( envname != "" ) {
		rv += GetEnv( envname );
	}
	return rv;
}

//----------------------------------------------------------------------------
// Notify system that things have changed
//----------------------------------------------------------------------------

void NotifyChanges() {
	const char * what = "Environment";
	SendMessageTimeout( HWND_BROADCAST, WM_SETTINGCHANGE, 0,
							(LPARAM) what, SMTO_ABORTIFHUNG, 50, NULL);
}


