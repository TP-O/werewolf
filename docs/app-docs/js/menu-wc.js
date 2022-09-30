'use strict';

customElements.define('compodoc-menu', class extends HTMLElement {
    constructor() {
        super();
        this.isNormalMode = this.getAttribute('mode') === 'normal';
    }

    connectedCallback() {
        this.render(this.isNormalMode);
    }

    render(isNormalMode) {
        let tp = lithtml.html(`
        <nav>
            <ul class="list">
                <li class="title">
                    <a href="index.html" data-type="index-link">Communication Server Template</a>
                </li>

                <li class="divider"></li>
                ${ isNormalMode ? `<div id="book-search-input" role="search"><input type="text" placeholder="Type to search"></div>` : '' }
                <li class="chapter">
                    <a data-type="chapter-link" href="index.html"><span class="icon ion-ios-home"></span>Getting started</a>
                    <ul class="links">
                        <li class="link">
                            <a href="overview.html" data-type="chapter-link">
                                <span class="icon ion-ios-keypad"></span>Overview
                            </a>
                        </li>
                        <li class="link">
                            <a href="index.html" data-type="chapter-link">
                                <span class="icon ion-ios-paper"></span>README
                            </a>
                        </li>
                        <li class="link">
                            <a href="license.html"  data-type="chapter-link">
                                <span class="icon ion-ios-paper"></span>LICENSE
                            </a>
                        </li>
                                <li class="link">
                                    <a href="dependencies.html" data-type="chapter-link">
                                        <span class="icon ion-ios-list"></span>Dependencies
                                    </a>
                                </li>
                                <li class="link">
                                    <a href="properties.html" data-type="chapter-link">
                                        <span class="icon ion-ios-apps"></span>Properties
                                    </a>
                                </li>
                    </ul>
                </li>
                    <li class="chapter modules">
                        <a data-type="chapter-link" href="modules.html">
                            <div class="menu-toggler linked" data-toggle="collapse" ${ isNormalMode ?
                                'data-target="#modules-links"' : 'data-target="#xs-modules-links"' }>
                                <span class="icon ion-ios-archive"></span>
                                <span class="link-name">Modules</span>
                                <span class="icon ion-ios-arrow-down"></span>
                            </div>
                        </a>
                        <ul class="links collapse " ${ isNormalMode ? 'id="modules-links"' : 'id="xs-modules-links"' }>
                            <li class="link">
                                <a href="modules/AppModule.html" data-type="entity-link" >AppModule</a>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-AppModule-533d5345d3740e94600162e9f2ed1eaef550a1c7a46a3040f8862773cd81c3faf1a58f50ee4efc70e3708ee5a0e808bac989be9fbca2e7156a58a9b90196862a"' : 'data-target="#xs-injectables-links-module-AppModule-533d5345d3740e94600162e9f2ed1eaef550a1c7a46a3040f8862773cd81c3faf1a58f50ee4efc70e3708ee5a0e808bac989be9fbca2e7156a58a9b90196862a"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-AppModule-533d5345d3740e94600162e9f2ed1eaef550a1c7a46a3040f8862773cd81c3faf1a58f50ee4efc70e3708ee5a0e808bac989be9fbca2e7156a58a9b90196862a"' :
                                        'id="xs-injectables-links-module-AppModule-533d5345d3740e94600162e9f2ed1eaef550a1c7a46a3040f8862773cd81c3faf1a58f50ee4efc70e3708ee5a0e808bac989be9fbca2e7156a58a9b90196862a"' }>
                                        <li class="link">
                                            <a href="injectables/CronService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >CronService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                            <li class="link">
                                <a href="modules/CommunicationModule.html" data-type="entity-link" >CommunicationModule</a>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-CommunicationModule-c28e879d4fcaca241dbb36db247841cb2733adaed829067140836a34295d1e98d5ead12a14568e6fdc217a1c25a37bea9de1a96d3414abac4e77db8f17e4dce7"' : 'data-target="#xs-injectables-links-module-CommunicationModule-c28e879d4fcaca241dbb36db247841cb2733adaed829067140836a34295d1e98d5ead12a14568e6fdc217a1c25a37bea9de1a96d3414abac4e77db8f17e4dce7"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-CommunicationModule-c28e879d4fcaca241dbb36db247841cb2733adaed829067140836a34295d1e98d5ead12a14568e6fdc217a1c25a37bea9de1a96d3414abac4e77db8f17e4dce7"' :
                                        'id="xs-injectables-links-module-CommunicationModule-c28e879d4fcaca241dbb36db247841cb2733adaed829067140836a34295d1e98d5ead12a14568e6fdc217a1c25a37bea9de1a96d3414abac4e77db8f17e4dce7"' }>
                                        <li class="link">
                                            <a href="injectables/AuthService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >AuthService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/CommunicationGateway.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >CommunicationGateway</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/CommunicationService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >CommunicationService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/PrismaService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >PrismaService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                            <li class="link">
                                <a href="modules/MessageModule.html" data-type="entity-link" >MessageModule</a>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-MessageModule-0b09c0b7d530d9da5b418b83669473c91b92106614faedc5d27ecb068fbe66152043874dc9630b1b234a8853428edc77afd9118441c8562bf75ea4e4c9938fda"' : 'data-target="#xs-injectables-links-module-MessageModule-0b09c0b7d530d9da5b418b83669473c91b92106614faedc5d27ecb068fbe66152043874dc9630b1b234a8853428edc77afd9118441c8562bf75ea4e4c9938fda"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-MessageModule-0b09c0b7d530d9da5b418b83669473c91b92106614faedc5d27ecb068fbe66152043874dc9630b1b234a8853428edc77afd9118441c8562bf75ea4e4c9938fda"' :
                                        'id="xs-injectables-links-module-MessageModule-0b09c0b7d530d9da5b418b83669473c91b92106614faedc5d27ecb068fbe66152043874dc9630b1b234a8853428edc77afd9118441c8562bf75ea4e4c9938fda"' }>
                                        <li class="link">
                                            <a href="injectables/MessageService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >MessageService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/PrismaService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >PrismaService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/RoomService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >RoomService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/UserService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >UserService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                            <li class="link">
                                <a href="modules/RoomModule.html" data-type="entity-link" >RoomModule</a>
                                    <li class="chapter inner">
                                        <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                            'data-target="#controllers-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' : 'data-target="#xs-controllers-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' }>
                                            <span class="icon ion-md-swap"></span>
                                            <span>Controllers</span>
                                            <span class="icon ion-ios-arrow-down"></span>
                                        </div>
                                        <ul class="links collapse" ${ isNormalMode ? 'id="controllers-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' :
                                            'id="xs-controllers-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' }>
                                            <li class="link">
                                                <a href="controllers/RoomController.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >RoomController</a>
                                            </li>
                                        </ul>
                                    </li>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' : 'data-target="#xs-injectables-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' :
                                        'id="xs-injectables-links-module-RoomModule-3701567b4070b501df3ce37d7a1496d1c7197bc52dd95ef574bfc452ddd67024e3867fa3f8a433203d81c1341388a7c80afcea5de228e14c06e1e32f74b78e61"' }>
                                        <li class="link">
                                            <a href="injectables/PrismaService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >PrismaService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/RoomService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >RoomService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                            <li class="link">
                                <a href="modules/UserModule.html" data-type="entity-link" >UserModule</a>
                                    <li class="chapter inner">
                                        <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                            'data-target="#controllers-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' : 'data-target="#xs-controllers-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' }>
                                            <span class="icon ion-md-swap"></span>
                                            <span>Controllers</span>
                                            <span class="icon ion-ios-arrow-down"></span>
                                        </div>
                                        <ul class="links collapse" ${ isNormalMode ? 'id="controllers-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' :
                                            'id="xs-controllers-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' }>
                                            <li class="link">
                                                <a href="controllers/UserController.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >UserController</a>
                                            </li>
                                        </ul>
                                    </li>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' : 'data-target="#xs-injectables-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' :
                                        'id="xs-injectables-links-module-UserModule-b1435eea3a9b66764cc2ea65aa8aeca9b6f2acb51a806970c6a377301654ef0fd7e7f9f473637e8fb068755cc186e23a5809a8cb01c7b0abc179ea894e29d7f7"' }>
                                        <li class="link">
                                            <a href="injectables/AuthService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >AuthService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/PrismaService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >PrismaService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/UserService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >UserService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                </ul>
                </li>
                    <li class="chapter">
                        <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ? 'data-target="#classes-links"' :
                            'data-target="#xs-classes-links"' }>
                            <span class="icon ion-ios-paper"></span>
                            <span>Classes</span>
                            <span class="icon ion-ios-arrow-down"></span>
                        </div>
                        <ul class="links collapse " ${ isNormalMode ? 'id="classes-links"' : 'id="xs-classes-links"' }>
                            <li class="link">
                                <a href="classes/AddToRoomDto.html" data-type="entity-link" >AddToRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/AllExceptionFilter.html" data-type="entity-link" >AllExceptionFilter</a>
                            </li>
                            <li class="link">
                                <a href="classes/BookRoomDto.html" data-type="entity-link" >BookRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/CreatePersistentRoomDto.html" data-type="entity-link" >CreatePersistentRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/CreatePersistentRoomsDto.html" data-type="entity-link" >CreatePersistentRoomsDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/CreateRoomDto.html" data-type="entity-link" >CreateRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/CreateTemporaryRoomsDto.html" data-type="entity-link" >CreateTemporaryRoomsDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/HttpExceptionFilter.html" data-type="entity-link" >HttpExceptionFilter</a>
                            </li>
                            <li class="link">
                                <a href="classes/InviteToRoomDto.html" data-type="entity-link" >InviteToRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/JoinRoomDto.html" data-type="entity-link" >JoinRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/KickOutOfRoomDto.html" data-type="entity-link" >KickOutOfRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/LeaveRoomDto.html" data-type="entity-link" >LeaveRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/MuteRoomDto.html" data-type="entity-link" >MuteRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/RedisIoAdapter.html" data-type="entity-link" >RedisIoAdapter</a>
                            </li>
                            <li class="link">
                                <a href="classes/RemoveFromRoomDto.html" data-type="entity-link" >RemoveFromRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/RemoveRoomsDto.html" data-type="entity-link" >RemoveRoomsDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/RespondRoomInvitationDto.html" data-type="entity-link" >RespondRoomInvitationDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/SendPrivateMessageDto.html" data-type="entity-link" >SendPrivateMessageDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/SendRoomMessageDto.html" data-type="entity-link" >SendRoomMessageDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/TransferOwnershipDto.html" data-type="entity-link" >TransferOwnershipDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/WsExceptionsFilter.html" data-type="entity-link" >WsExceptionsFilter</a>
                            </li>
                        </ul>
                    </li>
                        <li class="chapter">
                            <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ? 'data-target="#injectables-links"' :
                                'data-target="#xs-injectables-links"' }>
                                <span class="icon ion-md-arrow-round-down"></span>
                                <span>Injectables</span>
                                <span class="icon ion-ios-arrow-down"></span>
                            </div>
                            <ul class="links collapse " ${ isNormalMode ? 'id="injectables-links"' : 'id="xs-injectables-links"' }>
                                <li class="link">
                                    <a href="injectables/AuthService.html" data-type="entity-link" >AuthService</a>
                                </li>
                                <li class="link">
                                    <a href="injectables/EventNameBindingInterceptor.html" data-type="entity-link" >EventNameBindingInterceptor</a>
                                </li>
                                <li class="link">
                                    <a href="injectables/ParseIdPipe.html" data-type="entity-link" >ParseIdPipe</a>
                                </li>
                                <li class="link">
                                    <a href="injectables/PrismaService.html" data-type="entity-link" >PrismaService</a>
                                </li>
                                <li class="link">
                                    <a href="injectables/SocketUserIdBindingInterceptor.html" data-type="entity-link" >SocketUserIdBindingInterceptor</a>
                                </li>
                            </ul>
                        </li>
                    <li class="chapter">
                        <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ? 'data-target="#guards-links"' :
                            'data-target="#xs-guards-links"' }>
                            <span class="icon ion-ios-lock"></span>
                            <span>Guards</span>
                            <span class="icon ion-ios-arrow-down"></span>
                        </div>
                        <ul class="links collapse " ${ isNormalMode ? 'id="guards-links"' : 'id="xs-guards-links"' }>
                            <li class="link">
                                <a href="guards/AuthGuard.html" data-type="entity-link" >AuthGuard</a>
                            </li>
                        </ul>
                    </li>
                    <li class="chapter">
                        <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ? 'data-target="#miscellaneous-links"'
                            : 'data-target="#xs-miscellaneous-links"' }>
                            <span class="icon ion-ios-cube"></span>
                            <span>Miscellaneous</span>
                            <span class="icon ion-ios-arrow-down"></span>
                        </div>
                        <ul class="links collapse " ${ isNormalMode ? 'id="miscellaneous-links"' : 'id="xs-miscellaneous-links"' }>
                            <li class="link">
                                <a href="miscellaneous/enumerations.html" data-type="entity-link">Enums</a>
                            </li>
                            <li class="link">
                                <a href="miscellaneous/functions.html" data-type="entity-link">Functions</a>
                            </li>
                            <li class="link">
                                <a href="miscellaneous/typealiases.html" data-type="entity-link">Type aliases</a>
                            </li>
                            <li class="link">
                                <a href="miscellaneous/variables.html" data-type="entity-link">Variables</a>
                            </li>
                        </ul>
                    </li>
                    <li class="chapter">
                        <a data-type="chapter-link" href="coverage.html"><span class="icon ion-ios-stats"></span>Documentation coverage</a>
                    </li>
                    <li class="divider"></li>
                    <li class="copyright">
                        Documentation generated using <a href="https://compodoc.app/" target="_blank">
                            <img data-src="images/compodoc-vectorise.png" class="img-responsive" data-type="compodoc-logo">
                        </a>
                    </li>
            </ul>
        </nav>
        `);
        this.innerHTML = tp.strings;
    }
});