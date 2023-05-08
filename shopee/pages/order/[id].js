import { PayPalButtons, usePayPalScriptReducer } from "@paypal/react-paypal-js";
import axios from "axios";
import { useSession } from "next-auth/react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useReducer } from "react";
import { toast } from "react-toastify";
import { getError } from "@utils/error";
import { BaseLayout } from "@components/ui/layout";
import { useOrder } from "@components/hooks";
import { Loader } from "@components/ui/common";

function reducer(state, action) {
    switch (action.type) {
        case "PAY_REQUEST":
            return { ...state, loadingPay: true };
        case "PAY_SUCCESS":
            return { ...state, loadingPay: false, successPay: true };
        case "PAY_FAIL":
            return { ...state, loadingPay: false, errorPay: action.payload };
        case "PAY_RESET":
            return {
                ...state,
                loadingPay: false,
                successPay: false,
                errorPay: "",
            };

        case "DELIVER_REQUEST":
            return { ...state, loadingDeliver: true };
        case "DELIVER_SUCCESS":
            return { ...state, loadingDeliver: false, successDeliver: true };
        case "DELIVER_FAIL":
            return { ...state, loadingDeliver: false };
        case "DELIVER_RESET":
            return {
                ...state,
                loadingDeliver: false,
                successDeliver: false,
            };

        default:
            state;
    }
}
function OrderScreen() {
    const { data: session } = useSession();
    // order/:id
    const [{ isPending }, paypalDispatch] = usePayPalScriptReducer();

    const { query } = useRouter();
    const orderId = query.id;
    const { data: order } = useOrder(orderId)

    const [
        {
            successPay,
            loadingPay,
            loadingDeliver,
            successDeliver,
        },
        dispatch,
    ] = useReducer(reducer, {});

    useEffect(() => {
        if (successPay) {
            dispatch({ type: "PAY_RESET" });
        }
        if (successDeliver) {
            dispatch({ type: "DELIVER_RESET" });
        }
    }, [successDeliver, successPay]);

    useEffect(() => {
        if (order.data) {
            const loadPaypalScript = async () => {
                const { data: clientId } = await axios.get("/api/keys/paypal");
                paypalDispatch({
                    type: "resetOptions",
                    value: {
                        "client-id": clientId,
                        currency: "USD",
                    },
                });
                paypalDispatch({ type: "setLoadingStatus", value: "pending" });
            };
            loadPaypalScript();
        }
    }, [paypalDispatch, order.data])

    function createOrder(data, actions) {
        return actions.order
            .create({
                purchase_units: [
                    {
                        amount: { value: order.data.totalPrice },
                    },
                ],
            })
            .then((orderID) => {
                return orderID;
            });
    }

    function onApprove(data, actions) {
        return actions.order.capture().then(async function (details) {
            try {
                dispatch({ type: "PAY_REQUEST" });
                const { data } = await axios.put(
                    `/api/orders/${order.data._id}/pay`,
                    details
                );
                order.mutate({...order.data, 
                    isPaid: true
                })
                dispatch({ type: "PAY_SUCCESS", payload: data });
                toast.success("Order is paid successfully");
            } catch (err) {
                dispatch({ type: "PAY_FAIL", payload: getError(err) });
                toast.error(getError(err));
            }
        });
    }

    function onError(err) {
        toast.error(getError(err));
    }

    async function deliverOrderHandler() {
        try {
            dispatch({ type: "DELIVER_REQUEST" });
            const { data } = await axios.put(
                `/api/admin/orders/${order.data._id}/deliver`,
                {}
            );
            order.mutate({...order.data, 
                isDelivered: true
            })
            dispatch({ type: "DELIVER_SUCCESS", payload: data });
            toast.success("Order is delivered");
        } catch (err) {
            dispatch({ type: "DELIVER_FAIL", payload: getError(err) });
            toast.error(getError(err));
        }
    }

    return (
        <>
            <h1 className="mb-4 text-xl">{`Order ${orderId}`}</h1>
            {!order.hasInitialResponse ? (
                <div>Loading...</div>
            ) : order.error ? (
                <div className="alert-error">{order.error}</div>
            ) : (
                <div className="grid md:grid-cols-4 md:gap-5">
                    <div className="overflow-x-auto md:col-span-3">
                        <div className="card p-5">
                            <h2 className="mb-2 text-lg">Shipping Address</h2>
                            <div>
                                {order.data.shippingAddress.fullName},{" "}
                                {order.data.shippingAddress.address},{" "}
                                {order.data.shippingAddress.city},{" "}
                                {order.data.shippingAddress.postalCode},{" "}
                                {order.data.shippingAddress.country}
                            </div>
                            {order.data.isDelivered ? (
                                <div className="alert-success">
                                    Delivered at {order.data.deliveredAt}
                                </div>
                            ) : (
                                <div className="alert-error">Not delivered</div>
                            )}
                        </div>

                        <div className="card p-5">
                            <h2 className="mb-2 text-lg">Payment Method</h2>
                            <div>{order.data.paymentMethod}</div>
                            {order.data.isPaid ? (
                                <div className="alert-success">
                                    Paid at {order.data.paidAt}
                                </div>
                            ) : (
                                <div className="alert-error">Not paid</div>
                            )}
                        </div>

                        <div className="card overflow-x-auto p-5">
                            <h2 className="mb-2 text-lg">Order Items</h2>
                            <table className="min-w-full">
                                <thead className="border-b">
                                    <tr>
                                        <th className="px-5 text-left">Item</th>
                                        <th className="p-5 text-right">
                                            Quantity
                                        </th>
                                        <th className="p-5 text-right">
                                            Price
                                        </th>
                                        <th className="p-5 text-right">
                                            Subtotal
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {order.data.orderItems.map((item) => (
                                        <tr key={item._id} className="border-b">
                                            <td>
                                                <Link
                                                    href={`/product/${item.slug}`}
                                                >
                                                    <a className="flex items-center">
                                                        <Image
                                                            src={item.image}
                                                            alt={item.name}
                                                            width={50}
                                                            height={50}
                                                        ></Image>
                                                        &nbsp;
                                                        {item.name}
                                                    </a>
                                                </Link>
                                            </td>
                                            <td className="p-5 text-right">
                                                {item.quantity}
                                            </td>
                                            <td className="p-5 text-right">
                                                ${item.price}
                                            </td>
                                            <td className="p-5 text-right">
                                                ${item.quantity * item.price}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div>
                        <div className="card p-5">
                            <h2 className="mb-2 text-lg">Order Summary</h2>
                            <ul>
                                <li>
                                    <div className="mb-2 flex justify-between">
                                        <div>Items</div>
                                        <div>${order.data.itemsPrice}</div>
                                    </div>
                                </li>{" "}
                                <li>
                                    <div className="mb-2 flex justify-between">
                                        <div>Tax</div>
                                        <div>${order.data.taxPrice}</div>
                                    </div>
                                </li>
                                <li>
                                    <div className="mb-2 flex justify-between">
                                        <div>Shipping</div>
                                        <div>${order.data.shippingPrice}</div>
                                    </div>
                                </li>
                                <li>
                                    <div className="mb-2 flex justify-between">
                                        <div>Total</div>
                                        <div>${order.data.totalPrice}</div>
                                    </div>
                                </li>
                                {!order.data.isPaid && (
                                    <li>
                                        {isPending ? (
                                            <Loader/>
                                        ) : (
                                            <div className="w-full">
                                                <PayPalButtons
                                                    createOrder={createOrder}
                                                    onApprove={onApprove}
                                                    onError={onError}
                                                ></PayPalButtons>
                                            </div>
                                        )}
                                        {loadingPay && <Loader/>}
                                    </li>
                                )}
                                {session.user.isAdmin &&
                                    order.data.isPaid &&
                                    !order.data.isDelivered && (
                                        <li>
                                            {loadingDeliver && (
                                                <Loader/>
                                            )}
                                            <button
                                                className="primary-button w-full"
                                                onClick={deliverOrderHandler}
                                            >
                                                Deliver Order
                                            </button>
                                        </li>
                                    )}
                            </ul>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

OrderScreen.auth = true;
OrderScreen.Layout = BaseLayout;
export default OrderScreen;
