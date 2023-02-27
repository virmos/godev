import useSWRMutation from "swr/mutation";
import axios from "axios";

async function sendRequest(orderId, payload) {
    const { data } = await axios.put(
        `/api/orders/${orderId}/pay`,
        payload
    );
    return data;
}

export const handler = () => (orderId) => {
    const swrMutateRes = useSWRMutation(
        `/api/orders/${orderId}/pay`,
        sendRequest
    );

    return swrMutateRes;
};