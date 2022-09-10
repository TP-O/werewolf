import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common';
import { Reflector } from '@nestjs/core';

@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const roles = this.reflector.get<string[]>('roles', context.getHandler());

    console.log('Passed roles guard');

    if (!roles) {
      return true;
    }

    // const request = context.switchToHttp().getRequest();
    // const user = request.user;

    return true;
  }
}
