import { getSession } from 'next-auth/react';
import db from '@utils/db';

async function handler(req, res) {
    if (req.method !== 'PUT') {
        return res.status(400).send({ message: `${req.method} not supported` });
    }

    const session = await getSession({ req });
    if (!session) {
        return res.status(401).send({ message: 'signin required' });
    }

    const { user } = session;
    const { name, email, password } = req.body;

    if (
        !name ||
        !email ||
        !email.includes('@') ||
        (password && password.trim().length < 5)
    ) {
        res.status(422).json({
            message: 'Validation error',
        });
        return;
    }
    await db.connect();
    db.updateUserMongo(user, name, email, password)
    await db.disconnect();

    let payload = {
        _id: user._id,
        email: email,
        name: name,
        password: password
    }

    const requestOptions = {
        method: 'post',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    }
    
    return fetch("http://localhost:4000/api/update-user", requestOptions)
    .then(response => response.json())
    .then(user => {
        res.status(201).send({
            message: 'Updated user!',
            _id: user._id,
            name: user.name,
            email: user.email,
            isAdmin: user.is_admin,
        });
    }).catch(() => {
        throw new Error('Cannot update user!');
    })
}

export default handler;
