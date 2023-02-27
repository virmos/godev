import useSWRMutation from "swr/mutation";
import axios from "axios";

async function sendRequest(orderId, payload) {
    const { data } = await axios.put(
        `/api/admin/orders/${orderId}/deliver`,
        payload
    );
    return data;
}

export const handler = () => (orderId) => {
    const swrMutateRes = useSWRMutation(
        `/api/admin/orders/${orderId}/deliver`,
        sendRequest
    );

    return swrMutateRes;
};