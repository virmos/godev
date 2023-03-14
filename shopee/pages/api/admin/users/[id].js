import { getSession } from 'next-auth/react';

const handler = async (req, res) => {
    const session = await getSession({ req });
    if (!session || !session.user.isAdmin) {
        return res.status(401).send('admin signin required');
    }

    if (req.method === 'DELETE') {
        return deleteHandler(req, res);
    } else {
        return res.status(400).send({ message: 'Method not allowed' });
    }
};

const deleteHandler = async (req, res) => {
    let payload = {
        id: parseInt(req.query.id),
    }
    const requestOptions = {
        method: 'post',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    }
    
    return fetch("http://localhost:4000/api/get-user-by-id", requestOptions)
    .then(response => response.json())
    .then(user => {
        if (user?.id) {
            if (user.email === 'admin@example.com') {
                return res.status(400).send({ message: 'Can not delete admin' });
            }

            let payload = {
                id: parseInt(req.query.id),
            }

            const requestOptions = {
                method: 'post',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            }
            
            return fetch("http://localhost:4000/api/delete-user-by-id", requestOptions)
            .then(response => response.json())
            .then(() => {
                res.send({ message: 'User Deleted' });
            })

        } else {
            res.status(404).send({ message: 'User Not Found' });
        }
    }).catch(() => {
        throw new Error('Error deleting user');
    })
    
};

export default handler;
