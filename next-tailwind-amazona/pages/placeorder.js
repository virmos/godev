import Link from "next/link";
import { useRouter } from "next/router";
import React, { useEffect, useState } from "react";
import { toast } from "react-toastify";
import CheckoutWizard from "@components/ui/order/checkout/CheckoutWizard";
import { getError } from "@utils/error";
import { useStore } from "@components/providers";
import { BaseLayout } from "@components/ui/layout";
import { createOrder } from "@components/api";
import { Button, Message } from "@components/ui/common";
import {
    ShippingReport,
    PaymentReport,
    ProductReport,
    OrderReport,
} from "@components/ui/order";

export default function PlaceOrderScreen() {
    const { state, dispatch } = useStore();
    const { cart } = state;
    const { cartItems, shippingAddress, paymentMethod } = cart;

    const round2 = (num) => Math.round(num * 100 + Number.EPSILON) / 100;

    const itemsPrice = round2(
        cartItems.reduce((a, c) => a + c.quantity * c.price, 0)
    ); // 123.4567 => 123.46

    const shippingPrice = itemsPrice > 200 ? 0 : 15;
    const taxPrice = round2(itemsPrice * 0.15);
    const totalPrice = round2(itemsPrice + shippingPrice + taxPrice);

    const router = useRouter();
    useEffect(() => {
        if (!paymentMethod) {
            router.push("/payment");
        }
    }, [paymentMethod, router]);

    const [loading, setLoading] = useState(false);

    const placeOrderHandler = async () => {
        try {
            setLoading(true);
            const { data } = await createOrder({
                orderItems: cartItems,
                shippingAddress,
                paymentMethod,
                itemsPrice,
                shippingPrice,
                taxPrice,
                totalPrice,
            });
            setLoading(false);
            dispatch({ type: "CART_CLEAR_ITEMS" });
            router.push(`/order/${data._id}`);
        } catch (err) {
            setLoading(false);
            toast.error(getError(err));
        }
    };

    return (
        <>
            <CheckoutWizard activeStep={3} />
            <h1 className="mb-4 text-xl">Place Order</h1>
            {cartItems.length === 0 ? (
                <Message type="warning">
                    Cart is empty.
                    <Link href="/">
                        <a>
                            <Button size="sm" variant="purple" className="ml-2">
                                Go Shopping
                            </Button>
                        </a>
                    </Link>
                </Message>
            ) : (
                <div className="grid md:grid-cols-4 md:gap-5">
                    <div className="overflow-x-auto md:col-span-3">
                        <ShippingReport shippingAddress={shippingAddress} />
                        <PaymentReport paymentMethod={paymentMethod} />
                        <ProductReport cartItems={cartItems} />
                    </div>
                    <div>
                        <OrderReport
                            itemsPrice={itemsPrice}
                            taxPrice={taxPrice}
                            shippingPrice={shippingPrice}
                            totalPrice={totalPrice}
                            loading={loading}
                            placeOrderHandler={placeOrderHandler}
                        />
                    </div>
                </div>
            )}
        </>
    );
}

PlaceOrderScreen.auth = true;
PlaceOrderScreen.Layout = BaseLayout;
