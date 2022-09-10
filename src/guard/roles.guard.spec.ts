import { RolesGuard } from './roles.guard';

describe('RoleGuard', () => {
  it('should be defined', () => {
    expect(new RolesGuard(null)).toBeDefined();
  });
});
