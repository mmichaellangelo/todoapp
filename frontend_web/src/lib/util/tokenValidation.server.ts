import jwt from 'jsonwebtoken';

import type { SessionData } from '$lib/types';

export function getSessionDataFromToken(token: string): SessionData {
    const secretKey = "secret key"; // Should be in environment variable

    try {
        const decoded = jwt.verify(token, secretKey) as jwt.JwtPayload;
        
        if (!decoded || typeof decoded !== 'object') {
            throw new Error('Invalid token payload');
        }
        
        const username = decoded.username;
        const userid = decoded.userid as number;
        
        if (!username || typeof username !== 'string') {
            throw new Error('Username not found in token payload');
        }

        if (!userid) {
            throw new Error('User ID not found in token payload');
        }
        return {userid: userid, username: username};
    } catch (err) {
        throw err;
    }
}

