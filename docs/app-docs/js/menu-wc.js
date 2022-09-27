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
                            </li>
                            <li class="link">
                                <a href="modules/CommunicationModule.html" data-type="entity-link" >CommunicationModule</a>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-CommunicationModule-c2064feee3cf3036296ed14c4763f75c407fc858148cb35d402944ccc39368224f6f82540fc0449da745371d49e6e4a98eabc71ca2ba278d5e6e20544e0f469e"' : 'data-target="#xs-injectables-links-module-CommunicationModule-c2064feee3cf3036296ed14c4763f75c407fc858148cb35d402944ccc39368224f6f82540fc0449da745371d49e6e4a98eabc71ca2ba278d5e6e20544e0f469e"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-CommunicationModule-c2064feee3cf3036296ed14c4763f75c407fc858148cb35d402944ccc39368224f6f82540fc0449da745371d49e6e4a98eabc71ca2ba278d5e6e20544e0f469e"' :
                                        'id="xs-injectables-links-module-CommunicationModule-c2064feee3cf3036296ed14c4763f75c407fc858148cb35d402944ccc39368224f6f82540fc0449da745371d49e6e4a98eabc71ca2ba278d5e6e20544e0f469e"' }>
                                        <li class="link">
                                            <a href="injectables/AuthService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >AuthService</a>
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
                                        'data-target="#injectables-links-module-MessageModule-44210bb6323e15e7ca7d6db3af87c1c73fbc11ed2adde45740e3e4ed4aa5a481a6ecea31d75ebce263dd0aa5e5d5a0f89b4c3fd9365cbc63433e7c5560c51fa1"' : 'data-target="#xs-injectables-links-module-MessageModule-44210bb6323e15e7ca7d6db3af87c1c73fbc11ed2adde45740e3e4ed4aa5a481a6ecea31d75ebce263dd0aa5e5d5a0f89b4c3fd9365cbc63433e7c5560c51fa1"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-MessageModule-44210bb6323e15e7ca7d6db3af87c1c73fbc11ed2adde45740e3e4ed4aa5a481a6ecea31d75ebce263dd0aa5e5d5a0f89b4c3fd9365cbc63433e7c5560c51fa1"' :
                                        'id="xs-injectables-links-module-MessageModule-44210bb6323e15e7ca7d6db3af87c1c73fbc11ed2adde45740e3e4ed4aa5a481a6ecea31d75ebce263dd0aa5e5d5a0f89b4c3fd9365cbc63433e7c5560c51fa1"' }>
                                        <li class="link">
                                            <a href="injectables/MessageService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >MessageService</a>
                                        </li>
                                        <li class="link">
                                            <a href="injectables/PrismaService.html" data-type="entity-link" data-context="sub-entity" data-context-id="modules" >PrismaService</a>
                                        </li>
                                    </ul>
                                </li>
                            </li>
                            <li class="link">
                                <a href="modules/RoomModule.html" data-type="entity-link" >RoomModule</a>
                                <li class="chapter inner">
                                    <div class="simple menu-toggler" data-toggle="collapse" ${ isNormalMode ?
                                        'data-target="#injectables-links-module-RoomModule-32065ad5ecc159257cd716fab0e7aa2db3a7ff19e6b1274e34366968a5026a00f774d5ec4f5fec7ec2c58b288c36d5683b769ed94dfb119db8369cfac7cca135"' : 'data-target="#xs-injectables-links-module-RoomModule-32065ad5ecc159257cd716fab0e7aa2db3a7ff19e6b1274e34366968a5026a00f774d5ec4f5fec7ec2c58b288c36d5683b769ed94dfb119db8369cfac7cca135"' }>
                                        <span class="icon ion-md-arrow-round-down"></span>
                                        <span>Injectables</span>
                                        <span class="icon ion-ios-arrow-down"></span>
                                    </div>
                                    <ul class="links collapse" ${ isNormalMode ? 'id="injectables-links-module-RoomModule-32065ad5ecc159257cd716fab0e7aa2db3a7ff19e6b1274e34366968a5026a00f774d5ec4f5fec7ec2c58b288c36d5683b769ed94dfb119db8369cfac7cca135"' :
                                        'id="xs-injectables-links-module-RoomModule-32065ad5ecc159257cd716fab0e7aa2db3a7ff19e6b1274e34366968a5026a00f774d5ec4f5fec7ec2c58b288c36d5683b769ed94dfb119db8369cfac7cca135"' }>
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
                                <a href="classes/AllExceptionFilter.html" data-type="entity-link" >AllExceptionFilter</a>
                            </li>
                            <li class="link">
                                <a href="classes/BookRoomDto.html" data-type="entity-link" >BookRoomDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/CommunicationGateway.html" data-type="entity-link" >CommunicationGateway</a>
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
                                <a href="classes/RedisIoAdapter.html" data-type="entity-link" >RedisIoAdapter</a>
                            </li>
                            <li class="link">
                                <a href="classes/RespondRoomInvitationDto.html" data-type="entity-link" >RespondRoomInvitationDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/SendGroupMessageDto.html" data-type="entity-link" >SendGroupMessageDto</a>
                            </li>
                            <li class="link">
                                <a href="classes/SendPrivateMessageDto.html" data-type="entity-link" >SendPrivateMessageDto</a>
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