import { Module } from '@nestjs/common';
import { UserModule } from '../user/user.module';
import { RoomModule } from '../room/room.module';
import { AuthService } from 'src/common/service/auth.service';
import { PrismaService } from 'src/common/service/prisma.service';
import { CommunicationService } from './communication.service';
import { CommunicationGateway } from './communication.gateway';
import { FirebaseService } from 'src/common/service/firebase.service';

@Module({
  imports: [UserModule, RoomModule],
  providers: [
    CommunicationGateway,
    CommunicationService,
    AuthService,
    PrismaService,
    FirebaseService,
  ],
  exports: [CommunicationGateway],
})
export class CommunicationModule {}
