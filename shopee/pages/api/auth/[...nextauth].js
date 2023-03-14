// import bcryptjs from 'bcryptjs';
import NextAuth from 'next-auth';
import CredentialsProvider from 'next-auth/providers/credentials';
// import User from '@models/User';
// import db from '@utils/db';

export default NextAuth({
    session: {
        strategy: 'jwt',
    },
    callbacks: {
        async jwt({ token, user }) {
            if (user?._id) token._id = user._id;
            if (user?.isAdmin) token.isAdmin = user.isAdmin;
            return token;
        },
        async session({ session, token }) {
            if (token?._id) session.user._id = token._id;
            if (token?.isAdmin) session.user.isAdmin = token.isAdmin;
            return session;
        },
    },
    providers: [
        CredentialsProvider({
            async authorize(credentials) {
                let payload = {
                    email: credentials.email,
                    password: credentials.password
                }

                const requestOptions = {
                    method: 'post',
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                }
                return fetch("http://localhost:4000/api/login", requestOptions)
                .then(response => response.json())
                .then(user => {
                    return {
                        _id: user.id,
                        name: user.name,
                        email: user.email,
                        image: 'f',
                        isAdmin: user.is_admin,
                    };
                }).catch(() => {
                    throw new Error('Invalid email or password');
                })
            },
        }),
    ],
});
