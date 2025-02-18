import { NgModule, inject } from '@angular/core';
import { provideApollo } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { onError } from '@apollo/client/link/error';
import { WebSocketLink } from '@apollo/client/link/ws';
import { setContext } from '@apollo/client/link/context';
import { getMainDefinition } from '@apollo/client/utilities';
import { environment } from '../environments/environment.development';
import { InMemoryCache, split, ApolloLink } from '@apollo/client/core';

const httpUrl = environment.hasuraHttpEndpoint;
const wsUrl = environment.hasuraWsEndpoint;

interface Definition {
    kind: string;
    operation?: string;
}

@NgModule({
    providers: [
        provideApollo(() => {
            const httpLink = inject(HttpLink);
            const token = localStorage.getItem('access_token');

            const headers = setContext(() => {
                return {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                };
            });

            const errorLink = onError(({ graphQLErrors, networkError }) => {
                if (graphQLErrors) {
                    graphQLErrors.map(({ message, locations, path }) => {
                        console.error(
                            `GraphQL error\nMessage: ${message}\nLocation: ${locations}\nPath: ${path}`
                        );
                    });
                }
                if (networkError) {
                    console.log(`[Network error]: ${networkError}`);
                }
            });

            const httpClient = httpLink.create({
                uri: httpUrl,
            });

            const wsClient = new WebSocketLink({
                uri: wsUrl,
                options: {
                    reconnect: true,
                    connectionParams: async () => ({
                        headers: {
                            Authorization: `Bearer ${token ?? ''}`,
                        },
                    }),
                },
            });

            const link = errorLink.concat(
                split(
                    ({ query }) => {
                        const { kind, operation }: Definition = getMainDefinition(query);
                        return kind === 'OperationDefinition' && operation === 'subscription';
                    },
                    wsClient,
                    httpClient
                )
            );

            return {
                link: ApolloLink.from([headers, link]),
                cache: new InMemoryCache({ addTypename: false }),
                defaultOptions: {
                    watchQuery: { errorPolicy: 'all' },
                    mutate: { errorPolicy: 'all' },
                    query: { errorPolicy: 'all' },
                },
            };
        }),
    ],
})
export class GraphQLModule { }
