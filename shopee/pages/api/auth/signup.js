
async function handler(req, res) {
    if (req.method !== 'POST') {
        return;
    }
    const { name, email, password } = req.body;
    if (
        !name ||
        !email ||
        !email.includes('@') ||
        !password ||
        password.trim().length < 5
    ) {
        res.status(422).json({
            message: 'Validation error',
        });
        return;
    }

    let payload = {
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

    return fetch("http://localhost:4000/users/register", requestOptions)
    .then(response => response.json())
    .then(user => {
        res.status(201).send({
            message: 'Created user!',
            _id: user.id,
            name: user.name,
            email: user.email,
            isAdmin: user.is_admin,
        });
    }).catch(() => {
        throw new Error('Cannot create user!');
    })
}

export default handler;
