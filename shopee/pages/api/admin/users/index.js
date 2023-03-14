import { getSession } from 'next-auth/react';

const handler = async (req, res) => {
    const session = await getSession({ req });
    if (!session || !session.user.isAdmin) {
        return res.status(401).send('admin signin required');
    }
    let payload = {
    }

    const requestOptions = {
        method: 'post',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    }

    return fetch("http://localhost:4000/api/get-users", requestOptions)
    .then(response => response.json())
    .then(users => {
        res.status(201).send(users);
    }).catch(() => {
        throw new Error('Cannot fetch users!');
    })
};

export default handler;
