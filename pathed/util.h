//----------------------------------------------------------------------------
// util.h
//
// utility functions
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------


#ifndef INC_PATHED_UTIL_H
#define INC_PATHED_UTIL_H

#include <string>

std::string GetEnv( const std::string &  name );
std::string ExpandPath( const std::string & adir );
void NotifyChanges();


#endif
