import axios from "axios" 

export const getOrder = (orderId) => {
    return axios.get(`/api/orders/${orderId}`)
};

export const getProduct = (productId) => {
    return axios.get(`/api/products/${productId}`)
};

export const getAdminOrder = (orderId) => {
    return axios.get(`/api/admin/orders/${orderId}`);
};

export const getAdminProduct = (productId) => {
    return axios.get(`/api/admin/products/${productId}`);
};

export const getAdminUser = (userId) => {
    return axios.get(`/api/admin/users/${userId}`);
};

export const updateUserProfile = (newProfile) => {
    return axios.put('/api/auth/update', newProfile);
};

export const createNewProfile = (newProfile) => {
    return axios.post('/api/auth/signup', newProfile);
};

export const createOrder = (newOrder) => {
    return axios.post('/api/orders', newOrder);
};

export const deleteAdminProduct = (productId) => {
    return axios.delete(`/api/admin/products/${productId}`);
};

export const createAdminProducts = (newProduct) => {
    return axios.post(`/api/admin/products`, newProduct);
};

export const deleteAdminUser = (userId) => {
    return axios.delete(`/api/admin/users/${userId}`);
};

export const postOrder = (newOrder) => {
    return axios.post('/api/orders', newOrder);
};

export const updateAdminProduct = (productId, newProduct) => {
    return axios.put(`/api/admin/products/${productId}`, newProduct);
};
