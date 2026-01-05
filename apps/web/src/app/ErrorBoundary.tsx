import { Component, ErrorInfo, ReactNode } from 'react';
import { Button } from '../components/ui';

type ErrorBoundaryProps = {
    children: ReactNode;
};

type ErrorBoundaryState = {
    hasError: boolean;
    error: Error | null;
    errorInfo: ErrorInfo | null;
    detailsExpanded: boolean;
};

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
    state: ErrorBoundaryState = {
        hasError: false,
        error: null,
        errorInfo: null,
        detailsExpanded: false,
    };

    constructor(props: ErrorBoundaryProps) {
        super(props);
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return {
            hasError: true,
            error,
            errorInfo: null,
            detailsExpanded: false,
        };
    }

    componentDidCatch(error: Error, info: ErrorInfo): void {
        this.setState({ error, errorInfo: info });
    }

    resetError = () => {
        this.setState({
            hasError: false,
            error: null,
            errorInfo: null,
            detailsExpanded: false,
        });
    };

    render() {
        if (this.state.hasError) {
            return (
                <>
                    <main className='relative top-8 mb-8 flex h-fit w-11/12 justify-center rounded-lg bg-base-100 px-8 py-8 shadow xl:w-50% xl:min-w-50%'>
                        <div className='flex w-11/12 flex-col gap-4'>
                            <h1 className='text-2xl'>Error: {this.state.error && this.state.error.message}</h1>
                            <p className='flex flex-row items-center justify-between'>
                                There was an error with this page. Please return to the home page.
                                <Button
                                    variant='default'
                                    tabIndex={-1}
                                    onClick={() => {
                                        this.resetError();
                                        window.location.href = '/';
                                    }}>
                                    Home
                                </Button>
                            </p>
                        </div>
                    </main>
                </>
            );
        }
        return this.props.children;
    }
}
