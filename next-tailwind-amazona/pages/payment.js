import { useRouter } from 'next/router';
import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import CheckoutWizard from '@components/ui/order/checkout/CheckoutWizard';
import { BaseLayout } from '@components/ui/layout';
import { useStore } from '@components/providers';
import { Button } from '@components/ui/common';

export default function PaymentScreen() {
    const [selectedPaymentMethod, setSelectedPaymentMethod] = useState('');

    const { state, dispatch } = useStore();
    const { cart } = state;
    const { shippingAddress, paymentMethod } = cart;

    const router = useRouter();

    const submitHandler = (e) => {
        e.preventDefault();
        if (!selectedPaymentMethod) {
            return toast.error('Payment method is required');
        }
        dispatch({ type: 'SAVE_PAYMENT_METHOD', payload: selectedPaymentMethod });

        router.push('/placeorder');
    };
    useEffect(() => {
        if (!shippingAddress.address) {
            return router.push('/shipping');
        }
        setSelectedPaymentMethod(paymentMethod || '');
    }, [paymentMethod, router, shippingAddress.address]);

    return (
        <>
            <CheckoutWizard activeStep={2} />
            <form className="mx-auto max-w-screen-md" onSubmit={submitHandler}>
                <h1 className="mb-4 text-xl">Payment Method</h1>
                {['PayPal', 'Stripe', 'CashOnDelivery'].map((payment) => (
                    <div key={payment} className="mb-4">
                        <input
                            name="paymentMethod"
                            className="p-2 outline-none focus:ring-0"
                            id={payment}
                            type="radio"
                            checked={selectedPaymentMethod === payment}
                            onChange={() => setSelectedPaymentMethod(payment)}
                        />

                        <label className="p-2" htmlFor={payment}>
                            {payment}
                        </label>
                    </div>
                ))}
                <div className="mb-4 flex justify-between">
                    <button
                        onClick={() => router.push('/shipping')}
                        type="button"
                        className="default-button"
                    >
                        Back
                    </button>
                    <Button
                        variant='yellow'
                    >Next</Button>
                </div>
            </form>
        </>
    );
}

PaymentScreen.auth = true;
PaymentScreen.Layout = BaseLayout;
