import jwt from 'jsonwebtoken';

export function getUsernameFromAccessToken(token: string): string {
    const secretKey = "kinda secret key"; // Should be in environment variable

    try {
        const decoded = jwt.verify(token, secretKey) as jwt.JwtPayload;
        
        if (!decoded || typeof decoded !== 'object') {
            throw new Error('Invalid token payload');
        }
        
        const username = decoded.username;
        
        if (!username || typeof username !== 'string') {
            throw new Error('Username not found in token payload');
        }

        console.log("Username: ", username);
        return username;
    } catch (err) {
        console.error('Token verification failed:', err);
        throw err;
    }
}

